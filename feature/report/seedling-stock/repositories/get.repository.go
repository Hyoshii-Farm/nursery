package seedlingstock

import (
	"math"
	"sort"
	"time"

	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
)

const (
	pageLimit    = 10
	deadlineDays = 60 * 7 // seedling growth deadline in days +60 weeks
	dateFormat   = "2006-01-02"
)

// parseDate parses a date string using the standard date format.
func parseDate(date string) (time.Time, error) {
	return time.Parse(dateFormat, date)
}

const latestPlantingHistoryCTE = `
	JOIN (
		SELECT DISTINCT ON (location_id)
			location_id,
			variant_id,
			planting_date
		FROM "PlantingHistory"
		WHERE action = 'INITIAL'
		  AND deleted_at IS NULL
		  AND is_active = TRUE
		ORDER BY location_id, planting_date DESC, id DESC
	) ph ON ph.location_id = l.id
`

// sumQuantity sums the quantity of a specific action within a date range.
func (r *Repository) sumQuantity(
	action string,
	start time.Time,
	end time.Time,
) (int64, error) {
	var total int64
	err := r.db.Table(`"public"."SeedlingStock"`).
		Select("COALESCE(SUM(quantity), 0)").
		Where("action = ?", action).
		Where("datetime BETWEEN ? AND ?", start, end).
		Where("deleted_at IS NULL").
		Where("is_active = TRUE").
		Scan(&total).Error
	return total, err
}

func (r *Repository) getHighestLowestByAction(
	action string,
	start, end time.Time,
) (int, string, int, string, error) {

	type row struct {
		VariantName string
		Total       int
	}

	var rows []row

	err := r.db.Table(`"SeedlingStock" s`).
		Select(`
			v.name AS variant_name,
			SUM(s.quantity) AS total
		`).
		Joins(`JOIN "Variant" v ON v.id = s.variant_id`).
		Where("s.action = ?", action).
		Where("s.datetime BETWEEN ? AND ?", start, end).
		Where("s.deleted_at IS NULL").
		Where("s.is_active = TRUE").
		Group("v.name").
		Scan(&rows).Error

	if err != nil || len(rows) == 0 {
		return 0, "", 0, "", err
	}

	highestQty := rows[0].Total
	lowestQty := rows[0].Total
	highestName := rows[0].VariantName
	lowestName := rows[0].VariantName

	for _, r := range rows {
		if r.Total > highestQty {
			highestQty = r.Total
			highestName = r.VariantName
		}
		if r.Total < lowestQty {
			lowestQty = r.Total
			lowestName = r.VariantName
		}
	}

	return highestQty, highestName, lowestQty, lowestName, nil
}

// calcGap calculates the gap percentage between two values.
func calcGap(current, last int64) float64 {
	if last == 0 {
		return 0
	}
	raw := (float64(current-last) / float64(last)) * 100
	return math.Round(raw*100) / 100
}

// GetKPI retrieves KPI metrics for the report.
func (r *Repository) GetKPI(startDate, endDate string, variantIDs []uint, before bool) (model.KPI, error) {
	var kpi model.KPI

	// NOTE: variantIDs and before are reserved for future KPI filtering
	_ = variantIDs
	_ = before

	start, err := parseDate(startDate)
	if err != nil {
		return kpi, err
	}
	end, err := parseDate(endDate)
	if err != nil {
		return kpi, err
	}

	// KPI comparison uses fixed bi-weekly (14-day) periods
	duration := end.Sub(start)
	lastStart := start.Add(-duration)
	lastEnd := start

	type kpiTarget struct {
		action string
		target *model.KPIMetric
	}

	targets := []kpiTarget{
		{"ADD", &kpi.NewStock},
		{"TAKEN", &kpi.RemovedStock},
		{"DEAD", &kpi.DeadStock},
	}

	for _, t := range targets {
		current, err := r.sumQuantity(t.action, start, end)
		if err != nil {
			return kpi, err
		}

		last, err := r.sumQuantity(t.action, lastStart, lastEnd)
		if err != nil {
			return kpi, err
		}

		t.target.CurrentQuantity = int(current)
		t.target.LastPeriodQuantity = int(last)
		t.target.GapPercentage = calcGap(current, last)
		highestQty, highestName, lowestQty, lowestName, err :=
			r.getHighestLowestByAction(t.action, start, end)

		if err != nil {
			return kpi, err
		}

		t.target.HighestQuantity = highestQty
		t.target.HighestName = highestName
		t.target.LowestQuantity = lowestQty
		t.target.LowestName = lowestName
	}

	return kpi, nil
}

// GetSeedByVariant retrieves seed data grouped by variant.
// Only active variants (is_active = TRUE) are included.
func (r *Repository) GetSeedByVariant(endDate string, variantIDs []uint) ([]model.SeedByVariant, error) {
	var needs []struct {
		VariantID   uint
		VariantName string
		NeedQty     int
	}
	end, err := parseDate(endDate)
	if err != nil {
		return nil, err
	}
	// 1. Get need per variant (capacity-based)
	needQuery := r.db.Table(`"Location" l`).
		Select(`
			v.id   AS variant_id,
			v.name AS variant_name,
			SUM(l.capacity) AS need_qty
		`).
		Joins(latestPlantingHistoryCTE).
		Joins(`JOIN "Variant" v ON v.id = ph.variant_id`).
		Where("v.is_active = TRUE").
		Where("l.is_active = TRUE").
		Where("ph.planting_date + INTERVAL '52 weeks' <= ?", end).
		Group("v.id, v.name")

	if len(variantIDs) > 0 {
		needQuery = needQuery.Where("v.id IN ?", variantIDs)
	}

	if err := needQuery.Scan(&needs).Error; err != nil {
		return nil, err
	}

	// 2. Get available stock per variant
	available, err := r.GetAvailableSeed(endDate, variantIDs)
	if err != nil {
		return nil, err
	}

	availMap := make(map[string]int, len(available))
	for _, a := range available {
		availMap[a.VariantName] = a.AvailableQuantity
	}

	// 3. Merge need and available data
	result := make([]model.SeedByVariant, 0, len(needs))
	for _, n := range needs {
		avail := availMap[n.VariantName]
		result = append(result, model.SeedByVariant{
			VariantName:       n.VariantName,
			NeedQuantity:      n.NeedQty,
			AvailableQuantity: avail,
			GapQuantity:       avail - n.NeedQty,
		})
	}

	return result, nil
}

// GetSeedByLocation retrieves seed data grouped by location.
// Only active locations (is_active = TRUE) and active variants are included.
func (r *Repository) GetSeedByLocation(endDate string, variantIDs []uint, locationIDs []uint) ([]model.SeedByLocation, error) {
	var rows []struct {
		LocationName string
		VariantName  string
		NeedQty      int
		PlantingDate time.Time
	}
	end, err := parseDate(endDate)
	if err != nil {
		return nil, err
	}

	query := r.db.Table(`"Location" l`).
		Select(`
			l.name AS location_name,
			v.name AS variant_name,
			l.capacity AS need_qty,
			ph.planting_date AS planting_date
		`).
		Joins(latestPlantingHistoryCTE).
		Joins(`JOIN "Variant" v ON v.id = ph.variant_id`).
		Where("l.is_active = TRUE").
		Where("v.is_active = TRUE").
		Where("ph.planting_date + INTERVAL '52 weeks' <= ?", end).
		Order("ph.planting_date ASC")

	if len(variantIDs) > 0 {
		query = query.Where("v.id IN ?", variantIDs)
	}
	if len(locationIDs) > 0 {
		query = query.Where("l.id IN ?", locationIDs)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Get available stock per variant
	available, err := r.GetAvailableSeed(endDate, variantIDs)
	if err != nil {
		return nil, err
	}

	availMap := make(map[string]int, len(available))
	for _, a := range available {
		availMap[a.VariantName] = a.AvailableQuantity
	}

	// Merge and calculate deadline
	now := end
	result := make([]model.SeedByLocation, 0, len(rows))

	// Copy available into a remaining map
	remainingMap := make(map[string]int)
	for k, v := range availMap {
		remainingMap[k] = v
	}

	for _, row := range rows {
		remaining := remainingMap[row.VariantName]

		nextPlantingDate := row.PlantingDate.AddDate(0, 0, deadlineDays)
		daysLeft := int(nextPlantingDate.Sub(now).Hours() / 24)

		gap := remaining - row.NeedQty

		result = append(result, model.SeedByLocation{
			LocationName:      row.LocationName,
			NeedQuantity:      row.NeedQty,
			AvailableQuantity: remaining,
			GapQuantity:       gap,
			PlantingDate:      nextPlantingDate.Format("2006-01-02"),
			Deadline:          daysLeft,
		})

		// subtract this location's need from remaining stock
		remainingMap[row.VariantName] -= row.NeedQty
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Deadline < result[j].Deadline
	})

	return result, nil
}

// GetAvailableSeed retrieves available seed quantities by variant.
// Only active variants (is_active = TRUE) are included.
func (r *Repository) GetAvailableSeed(endDate string, variantIDs []uint) ([]model.AvailableSeed, error) {
	var result []model.AvailableSeed

	baseQuery := `
	WITH latest_initial AS (
		SELECT DISTINCT ON (variant_id)
			variant_id,
			datetime,
			id,
			quantity
		FROM "SeedlingStock"
		WHERE action = 'INITIAL'
		  AND deleted_at IS NULL
		  AND is_active = TRUE
		  AND datetime <= ?
		ORDER BY variant_id, datetime DESC, id DESC
	)
	SELECT 
		v.name AS variant_name,
		SUM(
			CASE
				WHEN s.action = 'INITIAL' THEN s.quantity
				WHEN s.action = 'ADD' THEN s.quantity
				WHEN s.action IN ('TAKEN','DEAD') THEN -s.quantity
				ELSE 0
			END
		) AS available_quantity
	FROM latest_initial li
	JOIN "SeedlingStock" s 
		ON s.variant_id = li.variant_id
		AND (
			s.datetime > li.datetime
			OR (s.datetime = li.datetime AND s.id >= li.id)
		)
	JOIN "Variant" v ON v.id = s.variant_id
	WHERE s.deleted_at IS NULL
	  AND s.is_active = TRUE
	  AND v.is_active = TRUE
	  AND s.datetime <= ?
	`

	args := []interface{}{endDate, endDate}

	if len(variantIDs) > 0 {
		baseQuery += " AND li.variant_id IN (?)"
		args = append(args, variantIDs)
	}

	baseQuery += " GROUP BY v.name"

	err := r.db.Raw(baseQuery, args...).Scan(&result).Error
	return result, err
}

// GetHistory retrieves paginated history records.
func (r *Repository) GetHistory(
	startDate, endDate string,
	variantIDs []uint,
	page uint,
) (model.History, error) {

	var history model.History
	records := []model.HistoryRecord{}

	if page == 0 {
		page = 1
	}
	offset := int((page - 1) * pageLimit)

	query := r.db.Table(`"public"."SeedlingStock" s`).
		Select(`
		TO_CHAR(s.datetime, 'DD FMMonth YYYY') AS date,
		s.action,
		v.name AS variant_name,
		s.quantity,
		s.description AS note
	`).
		Joins(`JOIN "Variant" v ON v.id = s.variant_id`).
		Where("s.datetime BETWEEN ? AND ?", startDate, endDate)

	if len(variantIDs) > 0 {
		query = query.Where("s.variant_id IN ?", variantIDs)
	}

	// Total count for pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return history, err
	}

	// Fetch paginated data
	if err := query.
		Order("s.datetime DESC").
		Limit(pageLimit).
		Offset(offset).
		Scan(&records).Error; err != nil {
		return history, err
	}

	pages := int(math.Ceil(float64(total) / float64(pageLimit)))

	history.Pagination = model.Pagination{
		Total: int(total),
		Page:  int(page),
		Limit: pageLimit,
		Pages: pages,
	}
	history.Data = records

	return history, nil
}
