package entities

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Province      Province `gorm:"not null;foreignKey:ProvinceID"`
	ProvinceID    uint     `gorm:"index"`
	City          City     `gorm:"not null;foreignKey:CityID"`
	CityID        uint     `gorm:"index"`
	StreetAddress string   `gorm:"not null"`
	HouseNumber   uint
	Unit          uint    `gorm:"default=1"`
	PostalCode    *string `gorm:"index"`
	OwnerID       uint    `gorm:"index"`
}
