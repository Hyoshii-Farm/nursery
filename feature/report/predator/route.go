package predator

import (
	handlers "github.com/Hyoshii-Farm/nursery/feature/report/predator/handlers"
	services "github.com/Hyoshii-Farm/nursery/feature/report/predator/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(router fiber.Router, db *gorm.DB) {
	service := services.NewService(db)
	handler := handlers.NewHandler(service)

	router.Get("/", handler.GetReport)
}
