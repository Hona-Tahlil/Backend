package request

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type CreateRequestRequest struct {
	UserID          uint
	PetSitterUserID uint
	CalenderSlots   []RequestCalendarSlotRequest
	PetIDs          []uint
	Notes           *string
	AddressID       uint
	ServiceIDs      []uint
}

type RequestCalendarSlotRequest struct {
	Date  time.Time
	Slots []enums.Slot
}
