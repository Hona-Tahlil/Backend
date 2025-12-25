package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	PetSitterID uint              `gorm:"index"`
	RequestID   *uint             `gorm:"index"`
	Type        enums.ServiceType `gorm:"index"`
	Price       uint              `gorm:"index"`
	Description *string           `gorm:"type:text"`
	Kind        string            `gorm:"index;default=petSitter"` // petSitter or request
}
