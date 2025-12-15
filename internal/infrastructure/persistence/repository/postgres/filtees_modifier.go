package postgres

import (
	"hona/backend/internal/application/dto/general"

	"gorm.io/gorm"
)



type FilterModifier struct {
	Filters []general.Filter
}
func NewFilterModifier(filters []general.Filter) FilterModifier {
	return FilterModifier{Filters: filters}
}


func (f FilterModifier) Apply(db *gorm.DB) *gorm.DB {
	for _, filter := range f.Filters {
		if filter.Op == "IN" {
			db = db.Where(filter.Field+" IN (?)", filter.Value)
			continue
		}
		db = db.Where(filter.Field+" "+filter.Op+" ?", filter.Value)
	}
	return db
}