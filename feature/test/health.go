package test

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
}

func HealthCheck(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "nursery-api",
		Version:   "0.0.1",
	}

	return c.JSON(response)
}
