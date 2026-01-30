package basic

type SortOption struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
