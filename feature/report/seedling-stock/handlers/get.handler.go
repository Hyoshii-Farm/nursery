package seedlingstock

import (
	"strconv"
	"strings"

	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
	"github.com/gofiber/fiber/v2"
)

// parseUintArray parses a comma-separated string of uints
func parseUintArray(param string) ([]uint, error) {
	if param == "" {
		return nil, nil
	}

	parts := strings.Split(param, ",")
	result := make([]uint, 0, len(parts))

	for _, p := range parts {
		v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(v))
	}

	return result, nil
}

func (h *Handler) GetReport(c *fiber.Ctx) error {
	// Parse query parameters
	req := model.SeedlingStockReportRequest{
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
	}

	// Parse variantID (optional, comma-separated list)
	if variantIDStr := c.Query("variantID"); variantIDStr != "" {
		variantIDs, err := parseUintArray(variantIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid variantID parameter",
			})
		}
		req.VariantID = variantIDs
	}

	// Parse page (defaults to 1)
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid page parameter",
			})
		}
		req.Page = uint(page)
	} else {
		req.Page = 1
	}

	// Parse locationID (optional, comma-separated list)
	if locationIDStr := c.Query("locationID"); locationIDStr != "" {
		locationIDs, err := parseUintArray(locationIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid locationID parameter",
			})
		}
		req.LocationID = locationIDs
	}

	// Parse before (optional boolean)
	req.Before = c.Query("before") == "true"

	// Call service
	result, err := h.service.GetReport(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
