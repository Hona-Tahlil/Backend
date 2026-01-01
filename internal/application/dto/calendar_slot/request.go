package calendarslot

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type CalendarSlotRequest struct {
	Date  time.Time
	Slots []enums.Slot
}

type GetCalendarSlotsRequest struct {
	UserID uint
}

type UpdateFreeCalendarSlotsRequest struct {
	UserID uint
	Add    []CalendarSlotRequest
	Remove []CalendarSlotRequest
}
