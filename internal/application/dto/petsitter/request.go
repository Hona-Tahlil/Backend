package petsitter

type Filter struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type Sort struct {
	Field string `json:"field"`
	Dir   string `json:"dir"`
}

type SearchPetSittersRequest struct {
	Offset  int
	Limit   int
	Filters []Filter
	Sorts   []Sort
}
