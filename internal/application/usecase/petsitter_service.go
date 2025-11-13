package usecase

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/servicedto"
)

type PetSitterService interface {
	GetPetSitterFreeSlotsResponse(id uint) ([]calendarslot.CalendarSlotInfoResponse, error)
	GetServicesResponse(id uint) ([]servicedto.ServiceInfoResponse, error)
}
