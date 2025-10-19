package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID      uint
	PetSitterID uint
	ReserveID   uint
	Title       string
	Description string
	Rating      int
}
