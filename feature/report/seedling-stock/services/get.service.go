package seedlingstock

import (
	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
)

func (s *Service) GetReport(req model.SeedlingStockReportRequest) (*model.SeedlingStockReportResponse, error) {
	// Get KPI data
	kpi, err := s.repo.GetKPI(req.StartDate, req.EndDate, req.VariantID, req.Before)
	if err != nil {
		return nil, err
	}

	// Get seed by variant
	seedByVariant, err := s.repo.GetSeedByVariant(req.StartDate, req.EndDate, req.VariantID)
	if err != nil {
		return nil, err
	}

	// Get seed by location
	seedByLocation, err := s.repo.GetSeedByLocation(req.StartDate, req.EndDate, req.VariantID, req.LocationID)
	if err != nil {
		return nil, err
	}

	// Get available seed
	availableSeed, err := s.repo.GetAvailableSeed(req.VariantID)
	if err != nil {
		return nil, err
	}

	if !req.Before {
		needMap := make(map[string]int)
		for _, v := range seedByVariant {
			needMap[v.VariantName] += v.NeedQuantity
		}

		for i := range availableSeed {
			availableSeed[i].AvailableQuantity -= needMap[availableSeed[i].VariantName]
		}
	}

	// Get history with pagination
	history, err := s.repo.GetHistory(req.StartDate, req.EndDate, req.VariantID, req.Page)
	if err != nil {
		return nil, err
	}

	return &model.SeedlingStockReportResponse{
		KPI:            kpi,
		SeedByVariant:  seedByVariant,
		SeedByLocation: seedByLocation,
		AvailableSeed:  availableSeed,
		History:        history,
	}, nil
}
