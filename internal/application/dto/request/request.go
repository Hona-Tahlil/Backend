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
	ServiceID       uint
}

type EditRequestRequest struct {
	RequestID     uint
	UserID        uint
	CalenderSlots []RequestCalendarSlotRequest
	PetIDs        []uint
	Notes         *string
	AddressInfo   *AddressInfoRequest
	AddressID     *uint
	ServiceID     uint
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

type GetCreateRequestInfoRequest struct {
	PetSitterUserID uint
	UserID          uint
}

type CancelRequestRequest struct {
	RequestID uint
	UserID    uint
}

type GetRequestFullDataRequest struct {
	RequestID uint
	UserID    uint
}
