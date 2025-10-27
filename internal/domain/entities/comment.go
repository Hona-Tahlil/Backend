package entities

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID      uint `gorm:"index;not null"`
	PetSitterID uint `gorm:"index;not null"`
	RequestID   uint `gorm:"index;not null"`
	Title       string
	Description string
	Rating      uint `gorm:"index;not null"`
}
