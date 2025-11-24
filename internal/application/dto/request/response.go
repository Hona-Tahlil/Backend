package request

import (
	"hona/backend/internal/application/dto/address"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/servicedto"
)

type CreateRequestInfoResponse struct {
	Services          []servicedto.ServiceInfoResponse
	Addresses         []address.AddressInfoResponse
	Pets              []string
	FreeCalendarSlots []calendarslot.CalendarSlotInfoResponse
}

type RequestFullDataResponse struct {
	RequestID       uint `json:"requestID"`
	PetSitterUserID uint `json:"petSitterUserID"`
	Service         servicedto.ServiceInfoResponse
	Pets            []pet.PetBasicDataResponse
	CalendarSlots   []calendarslot.CalendarSlotInfoResponse
	Notes           *string `json:"notes"`
	TotalPrice      uint    `json:"totalPrice"`
	// TODO: comment
	Address    address.AddressInfoResponse
	Status     string `json:"status"`
	TransferID *uint  `json:"transferID"`
}
