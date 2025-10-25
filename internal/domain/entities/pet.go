package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	UserID uint   `gorm:"index"`
	Name   string `gorm:"not null"`
	// TODO: join kind and species
	Kind      enums.PetKind `gorm:"index"`
	Species   string        // enum
	BirthDate *time.Time
	// TODO: calculate
	IsAdult bool
	Gender  *enums.Gender
	Weight  *uint
	// HealthRecord *string
	PictureLink *string
	AboutPet    *string
}
