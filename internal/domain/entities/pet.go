package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	UserID      uint          `gorm:"index"`
	Name        string        `gorm:"not null"`
	Kind        enums.PetKind `gorm:"index"`
	Species     string        // enum
	BirthDate   *time.Time
	IsAdult     bool
	Gender      *enums.Gender
	Weight      *uint
	PictureLink *string
	AboutPet    *string
}
