package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID          uint            `gorm:"index"`
	CertificateKeys []string        `gorm:"type:text[]"`
	FileKeys        []string        `gorm:"type:text[]"`
	Requests        []Request       `gorm:"foreignKey:PetSitterID"`
	Services        []Service       `gorm:"foreignKey:PetSitterID"`
	PetKinds        []enums.PetKind `gorm:"type:integer[]"`
	Schedule        []CalendarSlot  `gorm:"foreignKey:Refer"`
	Comments        []Comment       `gorm:"foreignKey:PetSitterID"`
	Bio             *string
	Status          enums.PetSitterStatus `gorm:"type:integer;index"`
	OnboardingStep  enums.OnboardingStep  `gorm:"default:1"`
}
