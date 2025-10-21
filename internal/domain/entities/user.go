package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"not null;unique;index"`
	IsVerified  bool   `gorm:"default=false;index"`
	Password    string `gorm:"not null"`
	Username    string `gorm:"unique;index;not null"`
	FirstName   string
	LastName    string
	Address     *Address `gorm:"foreignKey:OwnerID"`
	Phone       *string  `gorm:"index"`
	Gender      enums.Gender
	BirthDate   *time.Time
	PictureLink *string
	Wallet      Wallet    `gorm:"foreignKey:UserID;not null"`
	Requests    []Request `gorm:"foreignKey:UserID"`
	Role        Role      `gorm:"foreignKey:RoleID"`
	RoleID      uint
	Pets        []Pet      `gorm:"foreignKey:UserID"`
	PetSitter   *PetSitter `gorm:"foreignKey:UserID"`
	Comments    []Comment  `gorm:"foreignKey:UserID"`
}
