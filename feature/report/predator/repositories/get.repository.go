package predator

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"
)

func (r *Repository) GetReport(startDate, endDate string) ([]model.PredatorDTO, error) {
	var predators []model.PredatorDTO

	query := r.db.Model(&model.Predator{}).
		Select("DISTINCT \"Predator\".name, \"Predator\".id").
		Joins("JOIN \"CompanyPermit\" cp ON \"Predator\".id = cp.permit_id").
		Where("cp.start_date <= ? AND cp.end_date >= ?", startDate, endDate).
		Order("\"Predator\".name ASC")

	result := query.Find(&predators)

	return predators, result.Error
}
