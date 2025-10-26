package entities

import (
	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID           uint           `gorm:"index"`
	CertificateLinks []string       `gorm:"type:text[]"`
	IsVerified       bool           `gorm:"default=false"`
	Requests         []Request      `gorm:"foreignKey:PetSitterID"`
	Services         []Service      `gorm:"foreignKey:PetSitterID"`
	Schedule         []CalenderSlot `gorm:"foreignKey:Refer"`
	Comments         []Comment      `gorm:"foreignKey:PetSitterID"`
	Bio              *string
}
