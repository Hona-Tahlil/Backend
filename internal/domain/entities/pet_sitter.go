package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID          uint           `gorm:"index"`
	CertificateKeys []string       `gorm:"type:text[]"`
	IsVerified      bool           `gorm:"default=false"`
	Requests        []Request      `gorm:"foreignKey:PetSitterID"`
	Services        []Service      `gorm:"foreignKey:PetSitterID"`
	Schedule        []CalendarSlot `gorm:"foreignKey:Refer"`
	Comments        []Comment      `gorm:"foreignKey:PetSitterID"`
	Bio             *string
	Status          enums.PetSitterStatus `gorm:"type:varchar(20);default:'draft';index"`
	OnboardingStep  enums.OnboardingStep  `gorm:"default:1"`
}
