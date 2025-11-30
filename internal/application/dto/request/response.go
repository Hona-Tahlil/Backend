package request

import (
	"hona/backend/internal/application/dto/address"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/servicedto"
	"time"
)

type CreateRequestInfoResponse struct {
	Services          []servicedto.ServiceInfoResponse        `json:"services"`
	Addresses         []address.AddressInfoResponse           `json:"addresses"`
	Pets              []string                                `json:"pets"`
	FreeCalendarSlots []calendarslot.CalendarSlotInfoResponse `json:"freeCalendarSlots"`
}

type RequestFullDataResponse struct {
	RequestID       uint                                    `json:"requestID"`
	PetSitterUserID uint                                    `json:"petSitterUserID"`
	Service         servicedto.ServiceInfoResponse          `json:"service"`
	Pets            []pet.PetBasicDataResponse              `json:"pets"`
	CalendarSlots   []calendarslot.CalendarSlotInfoResponse `json:"calendarSlots"`
	Notes           *string                                 `json:"notes"`
	TotalPrice      uint                                    `json:"totalPrice"`
	// TODO: comment
	Address    address.AddressInfoResponse `json:"address"`
	Status     string                      `json:"status"`
	TransferID *uint                       `json:"transferID"`
	UpdatedAt  time.Time                   `json:"updatedAt"`
}
