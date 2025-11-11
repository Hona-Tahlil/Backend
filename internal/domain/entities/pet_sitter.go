package entities

import (
	"hona/backend/internal/domain/enums"

	"gorm.io/gorm"
)

type PetSitter struct {
	gorm.Model
	UserID          uint            `gorm:"index"`
	CertificateKeys []string        `gorm:"type:text[]"`
	IsVerified      bool            `gorm:"default=false"`
	Requests        []Request       `gorm:"foreignKey:PetSitterID"`
	Services        []Service       `gorm:"foreignKey:PetSitterID"`
	PetKinds        []enums.PetKind `gorm:"not null"`
	Schedule        []CalendarSlot  `gorm:"foreignKey:Refer"`
	Comments        []Comment       `gorm:"foreignKey:PetSitterID"`
	Bio             *string
}
