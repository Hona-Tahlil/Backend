package postgres

import (
	"hona/backend/internal/application/dto/general"

	"gorm.io/gorm"
)

type SortModifier struct {
	Sorts []general.Sort
}

func NewSortModifier(sorts []general.Sort) SortModifier {
	return SortModifier{Sorts: sorts}
}	

func (s SortModifier) Apply(db *gorm.DB) *gorm.DB {
	for _, sort := range s.Sorts {
		db = db.Order(sort.Field + " " + sort.Dir)
	}
	return db
}
