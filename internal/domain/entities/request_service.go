package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type RequestService struct {
	gorm.Model
	RequestID   uint
	Type        enums.ServiceType `gorm:"index"`
	Price       uint              `gorm:"index"`
	Description *string
}
