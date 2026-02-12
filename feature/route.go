package feature

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// feature
	predator "github.com/Hyoshii-Farm/nursery/feature/report/predator"
	seedlingstock "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock"
	"github.com/Hyoshii-Farm/nursery/feature/test"
)

func RegisterAll(api fiber.Router, db *gorm.DB, redisClient *redis.Client) {
	seedlingstock.Register(api.Group("/seedling-stock"), db, redisClient)
	predator.Register(api.Group("/predator"), db)

	// Health check endpoint
	api.Get("/health", test.HealthCheck)
}
