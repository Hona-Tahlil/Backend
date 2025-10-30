package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email           string    `gorm:"not null;unique;index"`
	IsEmailVerified bool      `gorm:"default=false;index"`
	Password        string    `gorm:"not null"`
	FirstName       string    `gorm:"not null"`
	LastName        string    `gorm:"not null"`
	Address         []Address `gorm:"foreignKey:OwnerID"`
	Phone           *string   `gorm:"index"`
	IsPhoneVerified bool      `gorm:"default=false;index"`
	Gender          enums.Gender
	BirthDate       *time.Time
	PictureLink     *string
	Wallet          Wallet     `gorm:"foreignKey:UserID;not null"`
	Requests        []Request  `gorm:"foreignKey:UserID"`
	Roles           []Role     `gorm:"many2many:user_roles"`
	Pets            []Pet      `gorm:"foreignKey:UserID"`
	PetSitter       *PetSitter `gorm:"foreignKey:UserID"`
	Comments        []Comment  `gorm:"foreignKey:UserID"`
}
