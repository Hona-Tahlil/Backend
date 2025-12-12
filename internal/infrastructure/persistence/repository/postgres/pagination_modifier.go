package postgres

import "gorm.io/gorm"

type PaginationModifier struct {
	Offset int
	Limit  int
}

func (p PaginationModifier) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Offset).Limit(p.Limit)
}
