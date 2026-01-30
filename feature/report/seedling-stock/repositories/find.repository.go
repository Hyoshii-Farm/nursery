package seedlingstock

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
)

func (r *Repository) FindSeedlingStock(id uint) (model.SeedlingStock, error) {
	var seedling_stock model.SeedlingStock
	result := r.db.First(&seedling_stock, id)
	return seedling_stock, result.Error
}

func (r *Repository) FindSeedlingStockByName(name string) (model.SeedlingStock, error) {
	var seedling_stock model.SeedlingStock
	result := r.db.Where("name = ?", name).First(&seedling_stock)
	return seedling_stock, result.Error
}
