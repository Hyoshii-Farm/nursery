package seedlingstock

import (
	service "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/services"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}
