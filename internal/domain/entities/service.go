package entities

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	PetSitterID uint
	Type        string // make enum later
	Price       int
}
