package seedlingstock

type SeedlingStockDTO struct {
	Name string `json:"name"`
}

// KPIMetric represents a single KPI metric with comparison data
type KPIMetric struct {
	CurrentQuantity    int     `json:"current_quantity"`
	LastPeriodQuantity int     `json:"last_period_quantity"`
	GapPercentage      float64 `json:"gap_percentage"`
	HighestQuantity    int     `json:"highest_quantity"`
	HighestName        string  `json:"highest_name"`
	LowestQuantity     int     `json:"lowest_quantity"`
	LowestName         string  `json:"lowest_name"`
}

// KPI contains all KPI metrics for the report
type KPI struct {
	NewStock     KPIMetric `json:"new_stock"`
	RemovedStock KPIMetric `json:"removed_stock"`
	DeadStock    KPIMetric `json:"dead_stock"`
}

// SeedByVariant represents seed data grouped by variant
type SeedByVariant struct {
	VariantName       string `json:"variant_name"`
	NeedQuantity      int    `json:"need_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
	GapQuantity       int    `json:"gap_quantity"`
}

// SeedByLocation represents seed data grouped by location
type SeedByLocation struct {
	LocationName      string `json:"location_name"`
	NeedQuantity      int    `json:"need_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
	GapQuantity       int    `json:"gap_quantity"`
	PlantingDate      string `json:"planting_date"`
	Deadline          int    `json:"deadline"`
}

// AvailableSeed represents available seed for a variant
type AvailableSeed struct {
	VariantName       string `json:"variant_name"`
	AvailableQuantity int    `json:"available_quantity"`
}

// HistoryRecord represents a single history entry
type HistoryRecord struct {
	Date        string `json:"date"`
	Action      string `json:"action"`
	VariantName string `json:"variant_name"`
	Quantity    int    `json:"quantity"`
	Note        string `json:"note"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// History contains paginated history data
type History struct {
	Pagination Pagination      `json:"pagination"`
	Data       []HistoryRecord `json:"data"`
}

// SeedlingStockReportResponse is the main response structure
type SeedlingStockReportResponse struct {
	KPI            KPI              `json:"KPI"`
	SeedByVariant  []SeedByVariant  `json:"seed_by_variant"`
	SeedByLocation []SeedByLocation `json:"seed_by_location"`
	AvailableSeed  []AvailableSeed  `json:"available_seed"`
	History        History          `json:"history"`
}

// SeedlingStockReportRequest contains query parameters
type SeedlingStockReportRequest struct {
	StartDate  string `query:"startDate"`
	EndDate    string `query:"endDate"`
	VariantID  []uint `query:"variantID"`
	Page       uint   `query:"page"`
	Before     bool   `query:"before"`
	LocationID []uint `query:"locationID"`
}
