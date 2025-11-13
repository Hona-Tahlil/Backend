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
	Pets              []pet.PetBasicDataResponse
	FreeCalendarSlots []calendarslot.CalendarSlotInfoResponse
}
