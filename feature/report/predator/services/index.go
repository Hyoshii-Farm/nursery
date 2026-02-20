package predator

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"
	repo "github.com/Hyoshii-Farm/nursery/feature/report/predator/repositories"

	"gorm.io/gorm"
)

type Service struct {
	repo *repo.Repository
}

func NewService(db *gorm.DB) *Service {
	repo := repo.GetRepository(db)
	return &Service{repo}
}

func calcPagination(total, page, limit int) model.Pagination {
	pages := 0
	if limit > 0 && total > 0 {
		pages = (total + limit - 1) / limit
	}
	return model.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}
}
