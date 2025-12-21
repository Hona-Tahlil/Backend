package domainpostgres

import "hona/backend/internal/application/dto/general"

type QueryModifier interface {
	Apply(query interface{}) interface{}
}

type FilteringOptions struct {
	Filters []general.Filter
}

type PaginationOptions struct {
	Limit  int
	Offset int
}

type SortingOptions struct {
	Sorts []general.Sort
}

type QueryOptions struct {
	Pagination *PaginationOptions
	Sorting    *SortingOptions
	Filters    *FilteringOptions
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

func (q *QueryOptions) WithPagination(limit, offset int) *QueryOptions {
	q.Pagination = &PaginationOptions{
		Limit:  limit,
		Offset: offset,
	}
	return q
}

func (q *QueryOptions) WithSorting(sorts []general.Sort) *QueryOptions {
	q.Sorting = &SortingOptions{
		Sorts: []general.Sort{},
	}
	for _, s := range sorts {
		q.Sorting.Sorts = append(q.Sorting.Sorts, s)
	}
	return q
}

func (q *QueryOptions) WithFilters(filters []general.Filter) *QueryOptions {
	q.Filters = &FilteringOptions{
		Filters: []general.Filter{},
	}
	for _, f := range filters {
		q.Filters.Filters = append(q.Filters.Filters, f)
	}
	return q
}

func (q *QueryOptions) HasPagination() bool {
	return q.Pagination != nil
}

func (q *QueryOptions) HasSorting() bool {
	return q.Sorting != nil
}

func (q *QueryOptions) HasFilters() bool {
	return q.Filters != nil
}