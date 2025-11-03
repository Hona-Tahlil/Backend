package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	PetSitterID uint              `gorm:"index"`
	Type        enums.ServiceType `gorm:"index"`
	Price       uint              `gorm:"index"`
	Description *string
}
