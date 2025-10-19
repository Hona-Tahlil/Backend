package entities

import (
	"time"

	"gorm.io/gorm"
)

type CalenderSlot struct {
	gorm.Model
	StartTime        time.Time
	EndTime          time.Time
	IsDailyRepeated  bool
	IsWeeklyRepeated bool
	Status           string // make enum later
	Refer            uint
}
