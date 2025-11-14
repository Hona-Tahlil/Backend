package service

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/domain/entities"
)

type CalendarSlotService struct {
}

func NewCalendarSlotService() *CalendarSlotService {
	return &CalendarSlotService{}
}

func (cs *CalendarSlotService) GetCalendarSlotsResponse(calendarSlots []entities.CalendarSlot) []calendarslot.CalendarSlotInfoResponse {
	r := make([]calendarslot.CalendarSlotInfoResponse, 0)

	for _, slot := range calendarSlots {
		r = append(r, calendarslot.CalendarSlotInfoResponse{
			ID:    slot.ID,
			Date:  slot.Date,
			Slots: slot.Slots,
		})
	}

	return r
}
