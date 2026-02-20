package feature

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	// feature
	predator "github.com/Hyoshii-Farm/nursery/feature/report/predator"
	seedlingstock "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock"
	"github.com/Hyoshii-Farm/nursery/feature/test"
)

func RegisterAll(api fiber.Router, db *gorm.DB) {
	seedlingstock.Register(api.Group("/seedling-stock"), db)
	predator.Register(api.Group("/report/predator-stock"), db)

	// Health check endpoint
	api.Get("/health", test.HealthCheck)
}
