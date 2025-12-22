package domainpostgres

import (
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type QueryOptions interface {
	WithPagination(limit, offset int) *postgres.QueryOptions
	WithSorting(sorts []general.Sort) *postgres.QueryOptions
	WithFilters(filters []general.Filter) *postgres.QueryOptions
	HasPagination() bool
	HasSorting() bool
	HasFilters() bool
}
