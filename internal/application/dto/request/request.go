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
	AddressInfo     *AddressInfoRequest
	AddressID       *uint
	ServiceIDs      []uint
}

type AddressInfoRequest struct {
	ProvinceName  enums.Province
	CityName      enums.City
	StreetAddress string
	HouseNumber   uint
	Unit          uint
	PostalCode    *string
}

type RequestCalendarSlotRequest struct {
	Date  time.Time
	Slots []enums.Slot
}
