package postgres

import "gorm.io/gorm"

type Filter struct {
	Field string
	Op    string
	Value interface{}
}

type FilterModifier struct {
	Filters []Filter
}

func (f FilterModifier) Apply(db *gorm.DB) *gorm.DB {
	for _, filter := range f.Filters {
		db = db.Where(filter.Field+" "+filter.Op+" ?", filter.Value)
	}
	return db
}
