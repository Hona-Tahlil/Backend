package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type CalendarSlot struct {
	gorm.Model
	Date        time.Time            `gorm:"not null"`
	Slots       Slots                `gorm:"type:integer[]"`
	Status      enums.CalendarStatus `gorm:"index"`
	PetSitterID *uint                `gorm:"index"`
	RequestID   *uint                `gorm:"index"`
}
