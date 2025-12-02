package calendarslot

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type CalendarSlotInfoResponse struct {
	ID     uint                 `json:"id"`
	Date   time.Time            `json:"date"`
	Slots  []enums.Slot         `json:"slots"`
	Status enums.CalendarStatus `json:"status"`
}
