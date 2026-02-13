package seedlingstock

import (
	handlers "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/handlers"
	services "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/services"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Register(router fiber.Router, db *gorm.DB, redisClient *redis.Client) {
	service := services.NewService(db, redisClient)
	handler := handlers.NewHandler(service)

	router.Get("/", handler.GetReport)
}
