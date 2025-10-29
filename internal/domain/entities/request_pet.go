package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type RequestPet struct {
	gorm.Model
	RequestID   uint
	Name        string        `gorm:"not null"`
	Kind        enums.PetKind `gorm:"index"`
	Species     enums.Species
	BirthDate   *time.Time
	IsAdult     bool
	Gender      *enums.Gender
	Weight      *uint
	PictureLink *string
	AboutPet    *string
}
