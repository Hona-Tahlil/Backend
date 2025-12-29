package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID          uint           `gorm:"index"`
	CertificateKeys []string       `gorm:"type:text[]"`
	FileKeys        []string       `gorm:"type:text[]"`
	Requests        []Request      `gorm:"foreignKey:PetSitterID"`
	Services        []Service      `gorm:"foreignKey:PetSitterID"`
	PetKinds        PetKinds       `gorm:"type:integer[]"`
	Schedule        []CalendarSlot `gorm:"foreignKey:PetSitterID"`
	Comments        []Comment      `gorm:"foreignKey:PetSitterID"`
	// Rating          uint           `gorm:"index;not null"`
	Bio            *string
	Status         enums.PetSitterStatus `gorm:"type:integer;index"`
	OnboardingStep enums.OnboardingStep  `gorm:"default:1"`
	PriceKey             *uint                 `gorm:"column:price_key;->"` // read-only, comes from SQL alias
}
