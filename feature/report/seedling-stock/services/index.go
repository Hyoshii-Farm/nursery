package seedlingstock

import (
	repo "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/repositories"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type Service struct {
	repo  *repo.Repository
	redis *redis.Client
}

func NewService(db *gorm.DB, redisClient *redis.Client) *Service {
	repo := repo.GetRepository(db)
	return &Service{repo: repo, redis: redisClient}
}
