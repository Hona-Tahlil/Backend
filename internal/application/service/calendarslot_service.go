package service

import (
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
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

func (cs *CalendarSlotService) GetFreeMap(calendarSlots []entities.CalendarSlot) map[string]map[interface{}]bool {
	availableSlots := make(map[string]map[interface{}]bool)

	for _, psSlot := range calendarSlots {
		if psSlot.Status == enums.Free {
			dateKey := psSlot.Date.Format("2006-01-02")

			if _, exists := availableSlots[dateKey]; !exists {
				availableSlots[dateKey] = make(map[interface{}]bool)
			}

			for _, slot := range psSlot.Slots {
				availableSlots[dateKey][slot] = true
			}
		}
	}

	return availableSlots
}
