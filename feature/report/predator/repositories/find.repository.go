package predator

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"
)

func (r *Repository) FindPredator(id uint) (model.Predator, error) {
	var predator model.Predator
	result := r.db.First(&predator, id)
	return predator, result.Error
}

func (r *Repository) FindPredatorByName(name string) (model.Predator, error) {
	var predator model.Predator
	result := r.db.Where("name = ?", name).First(&predator)
	return predator, result.Error
}
