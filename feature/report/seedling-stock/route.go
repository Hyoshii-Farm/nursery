package seedlingstock

import (
	handlers "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/handlers"
	services "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(router fiber.Router, db *gorm.DB) {
	service := services.NewService(db)
	handler := handlers.NewHandler(service)

	router.Get("/", handler.GetReport)
}
