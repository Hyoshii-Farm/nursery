package predator

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetReport(c *fiber.Ctx) error {
	req := new(model.PredatorPageRequest)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse query parameters",
		})
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	query, err := req.ToQuery()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result, err := h.service.GetReport(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
