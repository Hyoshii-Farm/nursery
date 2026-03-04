package seedlingstock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	model "github.com/Hyoshii-Farm/nursery/feature/report/seedling-stock/models"
	"github.com/redis/go-redis/v9"
)

const cacheTTL = 5 * time.Minute

// cacheGet is a generic helper that handles the Redis cache-aside pattern.
// It tries Redis first, falls back to the fetch function, then stores the result.
func cacheGet[T any](ctx context.Context, rdb *redis.Client, key string, ttl time.Duration, fetch func() (T, error)) (T, error) {
	// Try cache first
	if rdb != nil {
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			var cached T
			if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
				return cached, nil
			}
		}
	}

	// Fallback to DB
	data, err := fetch()
	if err != nil {
		var zero T
		return zero, err
	}

	// Store in cache
	if rdb != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("cache marshal error for key %s: %v", key, err)
		} else {
			rdb.Set(ctx, key, bytes, ttl)
		}
	}

	return data, nil
}

func (s *Service) GetReport(ctx context.Context, req model.SeedlingStockReportRequest) (*model.SeedlingStockReportResponse, error) {
	// Get KPI data
	kpi, err := s.getKPICached(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get seed by variant
	seedByVariant, err := s.getSeedByVariantCached(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get seed by location
	seedByLocation, err := s.getSeedByLocationCached(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get available seed
	availableSeed, err := s.getAvailableSeedCached(ctx, req)
	if err != nil {
		return nil, err
	}

	// When Before is false, subtract the need quantities from available seed
	// so the report reflects stock after fulfilling planned variant needs.
	if !req.Before {
		needMap := make(map[string]int)
		for _, v := range seedByVariant {
			needMap[v.VariantName] += v.NeedQuantity
		}

		for i := range availableSeed {
			availableSeed[i].AvailableQuantity -= needMap[availableSeed[i].VariantName]
		}
	}

	// History is fetched directly (not cached) because it's paginated
	// and users expect real-time data when browsing history pages.
	history, err := s.repo.GetHistory(req.StartDate, req.EndDate, req.VariantID, req.Page)
	if err != nil {
		return nil, err
	}

	return &model.SeedlingStockReportResponse{
		KPI:            *kpi,
		SeedByVariant:  seedByVariant,
		SeedByLocation: seedByLocation,
		AvailableSeed:  availableSeed,
		History:        history,
	}, nil
}

func (s *Service) getKPICached(ctx context.Context, req model.SeedlingStockReportRequest) (*model.KPI, error) {
	key := fmt.Sprintf(
		"seedling:kpi:%s:%s:%d:%t",
		req.StartDate,
		req.EndDate,
		req.VariantID,
		req.Before,
	)

	data, err := cacheGet(ctx, s.redis, key, cacheTTL, func() (model.KPI, error) {
		return s.repo.GetKPI(req.StartDate, req.EndDate, req.VariantID, req.Before)
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Service) getSeedByVariantCached(ctx context.Context, req model.SeedlingStockReportRequest) ([]model.SeedByVariant, error) {
	key := fmt.Sprintf(
		"seedling:seed_by_variant:%s:%s:%v",
		req.StartDate,
		req.EndDate,
		req.VariantID,
	)

	return cacheGet(ctx, s.redis, key, cacheTTL, func() ([]model.SeedByVariant, error) {
		return s.repo.GetSeedByVariant(req.EndDate, req.VariantID)
	})
}

func (s *Service) getSeedByLocationCached(ctx context.Context, req model.SeedlingStockReportRequest) ([]model.SeedByLocation, error) {
	key := fmt.Sprintf(
		"seedling:seed_by_location:%s:%d:%d",
		req.EndDate,
		req.VariantID,
		req.LocationID,
	)

	return cacheGet(ctx, s.redis, key, cacheTTL, func() ([]model.SeedByLocation, error) {
		return s.repo.GetSeedByLocation(req.EndDate, req.VariantID, req.LocationID)
	})
}

func (s *Service) getAvailableSeedCached(ctx context.Context, req model.SeedlingStockReportRequest) ([]model.AvailableSeed, error) {
	key := fmt.Sprintf(
		"seedling:available_seed:%s:%v",
		req.EndDate,
		req.VariantID,
	)

	return cacheGet(ctx, s.redis, key, cacheTTL, func() ([]model.AvailableSeed, error) {
		return s.repo.GetAvailableSeed(req.EndDate, req.VariantID)
	})
}
