package seedlingstock

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
)

func (r *Repository) GetReport(startDate, endDate string) ([]model.SeedlingStockDTO, error) {
	var seedling_stocks []model.SeedlingStockDTO

	query := r.db.Model(&model.SeedlingStock{}).
		Select("DISTINCT \"SeedlingStock\".name, \"SeedlingStock\".id").
		Joins("JOIN \"CompanyPermit\" cp ON \"SeedlingStock\".id = cp.permit_id").
		Where("cp.start_date <= ? AND cp.end_date >= ?", startDate, endDate).
		Order("\"SeedlingStock\".name ASC")

	result := query.Find(&seedling_stocks)

	return seedling_stocks, result.Error
}
