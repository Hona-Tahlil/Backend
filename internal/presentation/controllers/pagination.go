package controllers

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/general"
)

func GetOffsetLimit(page, pageSize int) (int, int) {
	if page <= 0 {
		page = bootstrap.Run().Constants.Pagination.DefaultPage
	}
	if pageSize <= 0 {
		pageSize = bootstrap.Run().Constants.Pagination.DefaultPageSize
	}

	return (page - 1) * pageSize, pageSize
}

// FilterParams represents a filter with Field, Op, and Value
type FilterParams struct {
	Field string
	Op    string
	Value any
}

// SortParams represents a sort with Field and Dir
type SortParams struct {
	Field string
	Dir   string
}

// ToFilters converts FilterParams to general.Filter DTO
func ToFilters(filters []FilterParams) []general.Filter {
	result := make([]general.Filter, len(filters))
	for i, f := range filters {
		result[i] = general.Filter{
			Field: f.Field,
			Op:    f.Op,
			Value: f.Value,
		}
	}
	return result
}

// ToSorts converts SortParams to general.Sort DTO
func ToSorts(sorts []SortParams) []general.Sort {
	result := make([]general.Sort, len(sorts))
	for i, s := range sorts {
		result[i] = general.Sort{
			Field: s.Field,
			Dir:   s.Dir,
		}
	}
	return result
}

type PaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

type PaginatedResponse[T any] struct {
	Data       []T             `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

func NewPaginatedResponse[T any](data []T, totalItems int64, offset, limit int) *PaginatedResponse[T] {
	pageSize := limit
	currentPage := int(offset/limit) + 1

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	}

	return &PaginatedResponse[T]{
		Data: data,
		Pagination: &PaginationMeta{
			CurrentPage: currentPage,
			PageSize:    pageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			HasNextPage: currentPage < totalPages,
			HasPrevPage: currentPage > 1,
		},
	}
}
