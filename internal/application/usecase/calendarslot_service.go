package usecase

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/domain/entities"
)

type CalendarSlotService interface {
	GetCalendarSlotsResponse(calendarSlots []entities.CalendarSlot) []calendarslot.CalendarSlotInfoResponse
}
