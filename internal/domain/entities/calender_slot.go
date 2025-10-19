package entities

import (
	"hona/backend/internal/domain/enums"
	"time"

	"gorm.io/gorm"
)

type CalenderSlot struct {
	gorm.Model
	StartTime        time.Time `gorm:"not null"`
	EndTime          time.Time `gorm:"not null"`
	IsDailyRepeated  bool      `gorm:"default=false"`
	IsWeeklyRepeated bool      `gorm:"default=false"`
	Status           enums.CalenderStatus
	Refer            uint
}
