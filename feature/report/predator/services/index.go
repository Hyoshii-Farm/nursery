package predator

import (
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
