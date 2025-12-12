package petsitter

import "hona/backend/internal/infrastructure/dsl"

type SearchPetSittersRequest struct {
	Offset  int
	Limit   int
	Filters []dsl.Filter
	Sorts   []dsl.Sort
}
