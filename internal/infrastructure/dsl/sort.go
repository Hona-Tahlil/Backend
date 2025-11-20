package dsl

import "gorm.io/gorm"

type Sort struct {
	Field string
	Dir   string // ASC یا DESC
}

type Sorts struct {
	items []Sort
}

func NewSorts() *Sorts {
	return &Sorts{items: make([]Sort, 0)}
}

func (s *Sorts) Asc(field string) *Sorts {
	s.items = append(s.items, Sort{Field: field, Dir: "ASC"})
	return s
}

func (s *Sorts) Desc(field string) *Sorts {
	s.items = append(s.items, Sort{Field: field, Dir: "DESC"})
	return s
}

func (s *Sorts) Apply(db *gorm.DB) *gorm.DB {
	for _, sort := range s.items {
		db = db.Order(sort.Field + " " + sort.Dir)
	}
	return db
}
