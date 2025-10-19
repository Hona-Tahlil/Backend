package entities

import (
	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID       uint
	Certificates []string `gorm:"type:text[]"`
	IsVerified   bool
	Bio          string
	Requests     []Request      `gorm:"foreignKey:PetSitterID"`
	Services     []Service      `gorm:"foreignKey:PetSitterID"`
	Schedule     []CalenderSlot `gorm:"foreignKey:Refer"`
	Comments     []Comment      `gorm:"foreignKey:PetSitterID"`
}
