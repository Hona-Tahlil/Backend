package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name       enums.City `gorm:"index;not null"`
	ProvinceID uint
}
