package petsitter

import (
	"hona/backend/bootstrap"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type PetSitterCalendarController struct {
	petSitterService usecase.PetSitterService
}

func NewPetSitterCalendarController(petSitterService usecase.PetSitterService) *PetSitterCalendarController {
	return &PetSitterCalendarController{
		petSitterService: petSitterService,
	}
}

func (pc *PetSitterCalendarController) GetCalendarSlots(ctx *gin.Context) {
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetCalendarSlots(calendarslot.GetCalendarSlotsRequest{UserID: userID})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterCalendarController) UpdateFreeCalendarSlots(ctx *gin.Context) {
	type CalendarSlotParams struct {
		Date  time.Time    `json:"date" binding:"required"`
		Slots []enums.Slot `json:"slots" binding:"required"`
	}
	type UpdateCalendarParams struct {
		Add    []CalendarSlotParams `json:"add" validate:"omitempty,dive"`
		Remove []CalendarSlotParams `json:"remove" validate:"omitempty,dive"`
	}

	params := controllers.Receive[UpdateCalendarParams](ctx)
	userID := controllers.GetID(ctx)

	add := make([]calendarslot.CalendarSlotRequest, len(params.Add))
	for i, slot := range params.Add {
		add[i] = calendarslot.CalendarSlotRequest{
			Date:  slot.Date,
			Slots: slot.Slots,
		}
	}
	remove := make([]calendarslot.CalendarSlotRequest, len(params.Remove))
	for i, slot := range params.Remove {
		remove[i] = calendarslot.CalendarSlotRequest{
			Date:  slot.Date,
			Slots: slot.Slots,
		}
	}

	err := pc.petSitterService.UpdateFreeCalendarSlots(calendarslot.UpdateFreeCalendarSlotsRequest{
		UserID: userID,
		Add:    add,
		Remove: remove,
	})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, nil)
}
