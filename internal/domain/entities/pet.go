package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	UserID     uint          `gorm:"index"`
	RequestID  *uint         `gorm:"index"`
	Name       string        `gorm:"not null"`
	Kind       enums.PetKind `gorm:"index"`
	Species    enums.Species
	BirthDate  *time.Time
	IsAdult    bool
	Gender     enums.PetGender
	Weight     *float32
	PictureKey *string
	AboutPet   *string
	Type       string `gorm:"default:'regular'"` // regular, request
}
