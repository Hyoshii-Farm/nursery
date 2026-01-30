package utils

import (
	"strings"

	types "github.com/Hyoshii-Farm/nursery/types"

	"gorm.io/gorm"
)

func ApplySorting(query *gorm.DB, sortOptions []types.SortOption) *gorm.DB {
	if len(sortOptions) == 0 {
		return query.Order("id asc")
	}

	for _, so := range sortOptions {
		order := strings.ToLower(so.Order)
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		query = query.Order(so.Field + " " + order)
	}

	return query
}

func ParseSortParam(sortParam string) []types.SortOption {
	if sortParam == "" {
		return nil
	}

	var sortOptions []types.SortOption
	sortPairs := strings.Split(sortParam, ",")

	for _, pair := range sortPairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 0 {
			continue
		}

		field := strings.TrimSpace(parts[0])
		if field == "" {
			continue
		}

		order := "asc"
		if len(parts) > 1 {
			order = strings.TrimSpace(parts[1])
		}

		sortOptions = append(sortOptions, types.SortOption{
			Field: field,
			Order: order,
		})
	}

	return sortOptions
}
