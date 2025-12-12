package postgres

import "gorm.io/gorm"

type Sort struct {
	Field string
	Dir   string
}
type SortModifier struct {
	Sorts []Sort
}

func (s SortModifier) Apply(db *gorm.DB) *gorm.DB {
	for _, sort := range s.Sorts {
		db = db.Order(sort.Field + " " + sort.Dir)
	}
	return db
}
