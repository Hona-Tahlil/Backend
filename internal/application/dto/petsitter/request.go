package petsitter

import (
	"time"
)

type GetPetSitterRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

type FirstSubmit struct {
	// Status      *string `json:"status,omitempty"`
	UserID    uint
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	// Gender      enums.Gender `json:"gender,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	City        string     `json:"city,omitempty"`
	Address     string     `json:"address,omitempty"`

	Bio *string `json:"bio,omitempty"`
}

type SecondSubmit struct {
}

// type User struct {
// 	gorm.Model
// 	Email           string    `gorm:"not null;unique;index"`
// 	IsEmailVerified bool      `gorm:"default=false;index"`
// 	Password        string    `gorm:"not null"`
// 	FirstName       string    `gorm:"not null"`
// 	LastName        string    `gorm:"not null"`
// 	Address         []Address `gorm:"foreignKey:OwnerID"`
// 	Phone           *string   `gorm:"index"`
// 	IsPhoneVerified bool      `gorm:"default=false;index"`
// 	Gender          enums.Gender
// 	BirthDate       *time.Time
// 	PictureLink     *string
// 	Wallet          Wallet     `gorm:"foreignKey:UserID;not null"`
// 	Requests        []Request  `gorm:"foreignKey:UserID"`
// 	Roles           []Role     `gorm:"many2many:user_roles"`
// 	Pets            []Pet      `gorm:"foreignKey:UserID"`
// 	PetSitter       *PetSitter `gorm:"foreignKey:UserID"`
// 	Comments        []Comment  `gorm:"foreignKey:UserID"`
// }
