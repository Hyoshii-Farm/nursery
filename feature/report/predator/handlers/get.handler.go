package predator

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetReport(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	result, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(result)
}
