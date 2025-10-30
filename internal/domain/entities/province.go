package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type Province struct {
	gorm.Model
	Name   enums.Province `gorm:"not null;index"`
	Cities []City         `gorm:"ProvinceID;not null"`
}
