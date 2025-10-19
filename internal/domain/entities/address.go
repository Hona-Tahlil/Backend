package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Province      Province `gorm:"not null;index;foreignKey:ProvinceID"`
	ProvinceID    uint
	City          City `gorm:"not null;index;foreignKey:CityID"`
	CityID        uint
	StreetAddress string `gorm:"not null"`
	HouseNumber   uint   `gorm:"not null"`
	Unit          uint   `gorm:"default=1"`
	PostalCode    string `gorm:"index"`
	OwnerID       uint
}
