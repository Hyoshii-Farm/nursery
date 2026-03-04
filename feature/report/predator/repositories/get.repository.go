package predator

import (
	"time"

	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"
)

func (r *Repository) GetHistory(startDate, endDate string, predatorIDs []uint, page, limit int) ([]model.PredatorHistory, error) {
	var history []model.PredatorHistory
	offset := (page - 1) * limit
	query := r.db.
		Where("datetime >= ? AND datetime <= ? AND is_active = true", startDate, endDate)

	if len(predatorIDs) > 0 {
		query = query.Where("predator_id IN ?", predatorIDs)
	}

	err := query.Order("datetime DESC").Offset(offset).Limit(limit).Find(&history).Error
	return history, err
}

func (r *Repository) CountHistory(startDate, endDate string, predatorIDs []uint) (int64, error) {
	var count int64
	query := r.db.Model(&model.PredatorHistory{}).
		Where("datetime >= ? AND datetime <= ? AND is_active = true", startDate, endDate)

	if len(predatorIDs) > 0 {
		query = query.Where("predator_id IN ?", predatorIDs)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *Repository) GetActionSummary(startDate, endDate string, predatorIDs []uint) (map[string]model.KPIBlock, error) {
	var results []struct {
		Action      string
		Quantity    float64
		VariantName string
	}

	query := r.db.Table("Predator").
		Select("ph.action, SUM(ph.quantity) as quantity, \"Predator\".name as variant_name").
		Joins("LEFT JOIN \"PredatorHistory\" ph ON \"Predator\".id = ph.predator_id AND ph.datetime >= ? AND ph.datetime <= ? AND ph.is_active = true", startDate, endDate)

	if len(predatorIDs) > 0 {
		query = query.Where("\"Predator\".id IN ?", predatorIDs)
	}

	err := query.Group("ph.action, \"Predator\".id, \"Predator\".name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]model.KPIBlock)
	for _, res := range results {
		block, ok := summary[res.Action]
		if !ok {
			block = model.KPIBlock{
				LowestQuantity: res.Quantity,
				LowestName:     res.VariantName,
			}
		}

		block.CurrentQuantity += res.Quantity

		if res.Quantity > block.HighestQuantity {
			block.HighestQuantity = res.Quantity
			block.HighestName = res.VariantName
		}
		if res.Quantity <= block.LowestQuantity {
			block.LowestQuantity = res.Quantity
			block.LowestName = res.VariantName
		}
		summary[res.Action] = block
	}

	return summary, nil
}

// GetVariantSummary returns all variant quantities/ages at a date — used for KPI aggregation (no pagination).
func (r *Repository) GetVariantSummary(date string, predatorIDs []uint) (map[string]model.VariantKPI, error) {
	results, err := r.queryVariantRows(date, predatorIDs, 0, 0)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	variants := make(map[string]model.VariantKPI)
	for _, res := range results {
		oldest, youngest := calcAge(now, res.OldestDate, res.YoungestDate)
		variants[res.Name] = model.VariantKPI{
			VariantName: res.Name,
			Quantity:    res.Quantity,
			Oldest:      oldest,
			Youngest:    youngest,
		}
	}
	return variants, nil
}

// GetStockPaginated returns a paginated list of stock items at a specific date.
func (r *Repository) GetStockPaginated(date string, predatorIDs []uint, page, limit int) ([]model.StockItem, int64, error) {
	// Count total rows
	var count int64
	countQuery := r.db.Table("Predator").
		Select("COUNT(DISTINCT \"Predator\".id)").
		Joins("LEFT JOIN \"PredatorHistory\" ph ON \"Predator\".id = ph.predator_id AND ph.datetime <= ? AND ph.is_active = true", date)
	if len(predatorIDs) > 0 {
		countQuery = countQuery.Where("\"Predator\".id IN ?", predatorIDs)
	}
	if err := countQuery.Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	results, err := r.queryVariantRows(date, predatorIDs, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	items := make([]model.StockItem, 0, len(results))
	for _, res := range results {
		oldest, _ := calcAge(now, res.OldestDate, res.YoungestDate)
		items = append(items, model.StockItem{
			VariantName: res.Name,
			Quantity:    res.Quantity,
			Age:         oldest,
		})
	}
	return items, count, nil
}

type variantRow struct {
	Name         string
	Quantity     float64
	OldestDate   time.Time
	YoungestDate time.Time
}

func (r *Repository) queryVariantRows(date string, predatorIDs []uint, offset, limit int) ([]variantRow, error) {
	var results []variantRow
	query := r.db.Table("Predator").
		Select("\"Predator\".name, COALESCE(SUM(CASE WHEN ph.action = 'ADD' THEN ph.quantity WHEN ph.action = 'REMOVED' THEN -ph.quantity ELSE 0 END), 0) as quantity, MIN(CASE WHEN ph.action = 'ADD' THEN ph.datetime END) as oldest_date, MAX(CASE WHEN ph.action = 'ADD' THEN ph.datetime END) as youngest_date").
		Joins("LEFT JOIN \"PredatorHistory\" ph ON \"Predator\".id = ph.predator_id AND ph.datetime <= ? AND ph.is_active = true", date)

	if len(predatorIDs) > 0 {
		query = query.Where("\"Predator\".id IN ?", predatorIDs)
	}

	query = query.Group("\"Predator\".id, \"Predator\".name").Order("\"Predator\".name")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	return results, query.Scan(&results).Error
}

func calcAge(now time.Time, oldest, youngest time.Time) (float64, float64) {
	o, y := 0.0, 0.0
	if !oldest.IsZero() {
		o = now.Sub(oldest).Hours() / 24
	}
	if !youngest.IsZero() {
		y = now.Sub(youngest).Hours() / 24
	}
	return o, y
}
