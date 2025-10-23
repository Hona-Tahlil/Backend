package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"not null"`
	Password    string `gorm:"not null"`
	IsVerified  bool
	Name        string
	Gender      string // make enum later
	Address     string
	BirthDate   time.Time
	Phone       string
	City        string
	PictureLink string
	// Wallet      Wallet    `gorm:"foreignKey:UserID"`
	// Requests    []Request `gorm:"foreignKey:UserID"`
	Roles []Role `gorm:"many2many:user_role"`
	// Pets        []Pet     `gorm:"foreignKey:UserID"`
	// PetSitter   PetSitter `gorm:"foreignKey:UserID"`
	// Comments    []Comment `gorm:"foreignKey:UserID"`
}
