package request

import "time"

type CreateRequestRequest struct {
	UserID        uint
	PetSitterID   uint
	CalenderSlots RequestCalendarSlotRequest
	PetIDs        []uint
	Notes         *string
	AddressID     uint
	ServiceIDs    []uint
}

type RequestCalendarSlotRequest struct {
	StartTime        time.Time
	EndTime          time.Time
	IsDailyRepeated  bool
	IsWeeklyRepeated bool
}
