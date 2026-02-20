package predator

import (
	"fmt"
	"strings"
	"time"
)

type PredatorPageRequest struct {
	StartDate   string `query:"startDate"`
	EndDate     string `query:"endDate"`
	PredatorIDs string `query:"predatorID"` // Changed to string to handle comma-separated values manually if needed, or better yet let's keep it consistent
	PageStock   int    `query:"pageStock"`
	PageHistory int    `query:"pageHistory"`
	Limit       int    `query:"limit"`
}

type PredatorPageQuery struct {
	StartDate   time.Time
	EndDate     time.Time
	PredatorIDs []uint
	PageStock   int
	PageHistory int
	Limit       int
}

func (r *PredatorPageRequest) Validate() error {
	if r.StartDate == "" || r.EndDate == "" {
		return fmt.Errorf("startDate and endDate are required")
	}

	layout := "2006-01-02"
	start, err := time.Parse(layout, r.StartDate)
	if err != nil {
		return fmt.Errorf("invalid startDate format, expected yyyy-mm-dd")
	}

	end, err := time.Parse(layout, r.EndDate)
	if err != nil {
		return fmt.Errorf("invalid endDate format, expected yyyy-mm-dd")
	}

	if start.After(end) {
		return fmt.Errorf("startDate cannot be after endDate")
	}

	if r.PredatorIDs == "" {
		return fmt.Errorf("at least one predatorID is required")
	}

	if r.PageStock < 1 {
		r.PageStock = 1
	}

	if r.PageHistory < 1 {
		r.PageHistory = 1
	}

	if r.Limit < 1 {
		r.Limit = 10
	}

	return nil
}

func (r *PredatorPageRequest) ToQuery() (*PredatorPageQuery, error) {
	layout := "2006-01-02"

	start, err := time.Parse(layout, r.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid startDate format, expected yyyy-mm-dd")
	}

	end, err := time.Parse(layout, r.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid endDate format, expected yyyy-mm-dd")
	}

	if start.After(end) {
		return nil, fmt.Errorf("startDate cannot be after endDate")
	}

	// make end inclusive
	end = end.Add(24*time.Hour - time.Nanosecond)

	// Parse PredatorIDs from comma-separated string
	var predatorIDs []uint
	if r.PredatorIDs != "" {
		parts := strings.Split(r.PredatorIDs, ",")
		for _, p := range parts {
			var id uint
			fmt.Sscanf(p, "%d", &id)
			if id > 0 {
				predatorIDs = append(predatorIDs, id)
			}
		}
	}

	if r.PageStock < 1 {
		r.PageStock = 1
	}
	if r.PageHistory < 1 {
		r.PageHistory = 1
	}
	if r.Limit < 1 {
		r.Limit = 10
	}

	return &PredatorPageQuery{
		StartDate:   start,
		EndDate:     end,
		PredatorIDs: predatorIDs,
		PageStock:   r.PageStock,
		PageHistory: r.PageHistory,
		Limit:       r.Limit,
	}, nil
}

type PredatorPageResponse struct {
	KPI     KPISection     `json:"KPI"`
	Stock   StockSection   `json:"stock"`
	History HistorySection `json:"history"`
}

type KPISection struct {
	TotalStock   KPIBlock     `json:"total_stock"`
	NewStock     KPIBlock     `json:"new_stock"`
	RemovedStock KPIBlock     `json:"removed_stock"`
	AverageAge   AgeBlock     `json:"average_age"`
	Variant      []VariantKPI `json:"variant"`
}

type KPIBlock struct {
	CurrentQuantity    float64 `json:"current_quantity"`
	LastPeriodQuantity float64 `json:"last_period_quantity"`
	GapPercentage      float64 `json:"gap_percentage"`
	HighestQuantity    float64 `json:"highest_quantity"`
	HighestName        string  `json:"highest_name"`
	LowestQuantity     float64 `json:"lowest_quantity"`
	LowestName         string  `json:"lowest_name"`
}

type AgeBlock struct {
	CurrentAge    float64 `json:"current_age"`
	LastPeriodAge float64 `json:"last_period_age"`
	GapPercentage float64 `json:"gap_percentage"`
	Oldest        float64 `json:"oldest"`
	OldestName    string  `json:"oldest_name"`
	Youngest      float64 `json:"youngest"`
	YoungestName  string  `json:"youngest_name"`
}

type VariantKPI struct {
	VariantName        string  `json:"variant_name"`
	Quantity           float64 `json:"quantity"`
	LastPeriodQuantity float64 `json:"last_period_quantity"`
	GapQuantity        float64 `json:"gap_quantity"`
	Oldest             float64 `json:"oldest"`
	Youngest           float64 `json:"youngest"`
}

type StockSection struct {
	Pagination Pagination  `json:"pagination"`
	Data       []StockItem `json:"data"`
}

type StockItem struct {
	VariantName string  `json:"variant_name"`
	Quantity    float64 `json:"quantity"`
	Age         float64 `json:"age"`
}

type HistorySection struct {
	Pagination Pagination    `json:"pagination"`
	Data       []HistoryItem `json:"data"`
}

type HistoryItem struct {
	Date        time.Time `json:"date"`
	Action      string    `json:"action"`
	VariantName string    `json:"variant_name"`
	Quantity    float64   `json:"quantity"`
	Note        string    `json:"note"`
}

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}
