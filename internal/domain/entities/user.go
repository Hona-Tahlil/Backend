package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"not null;unique;index"`
	Password    string `gorm:"not null"`
	Name        string `gorm:"not null;index"`
	Gender      enums.Gender
	Address     string
	BirthDate   time.Time
	Phone       string
	City        string
	PictureLink string
	Wallet      Wallet    `gorm:"foreignKey:UserID"`
	Requests    []Request `gorm:"foreignKey:UserID"`
	Roles       []Role    `gorm:"many2many:user_role"`
	Pets        []Pet     `gorm:"foreignKey:UserID"`
	PetSitter   PetSitter `gorm:"foreignKey:UserID"`
	Comments    []Comment `gorm:"foreignKey:UserID"`
}
