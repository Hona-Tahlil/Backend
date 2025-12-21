package postgres

import "gorm.io/gorm"

type QueryModifier interface {
	Apply(db *gorm.DB) *gorm.DB
}

