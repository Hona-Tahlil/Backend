package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	UserID       uint          `gorm:"index"`
	Name         string        `gorm:"not null"`
	Kind         enums.PetKind `gorm:"index"`
	Species      *string
	BirthDate    *time.Time
	Gender       enums.Gender
	Weight       *uint
	HealthRecord *string
	PictureLink  *string
	AboutPet     *string
}
