package user

import (
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/service"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type UserRequestController struct {
	requestService *service.RequestService
}

func NewUserRequestController(requestService *service.RequestService) *UserRequestController {
	return &UserRequestController{
		requestService: requestService,
	}
}

// TODO: Initialize Request / Rbac (Verified Email) / Validation / Email Sending? / chat / Status / Notification?
func (rc *UserRequestController) CreateRequest(ctx *gin.Context) {
	type CalendarSlot struct {
		Date  time.Time    `json:"date" validate:"required"`
		Slots []enums.Slot `json:"slots" validate:"required"`
	}
	type Params struct {
		PetSitterUserID uint           `json:"petSitterUserID" validate:"required"`
		CalenderSlots   []CalendarSlot `json:"calendarSlots" validate:"required"`
		PetIDs          []uint         `json:"petIDs" validate:"required"`
		Notes           *string        `json:"notes"`
		// TODO: make address
		AddressID  uint   `json:"addressID" validate:"required"`
		ServiceIDs []uint `json:"serviceIDs" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)
	slots := make([]request.RequestCalendarSlotRequest, 0)
	for _, slot := range params.CalenderSlots {
		slots = append(slots, request.RequestCalendarSlotRequest{
			Date:  slot.Date,
			Slots: slot.Slots,
		})
	}
	info := request.CreateRequestRequest{
		UserID:          UserID,
		PetSitterUserID: params.PetSitterUserID,
		CalenderSlots:   slots,
		PetIDs:          params.PetIDs,
		Notes:           params.Notes,
		AddressID:       params.AddressID,
		ServiceIDs:      params.ServiceIDs,
	}

	if err := rc.requestService.CreateRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

// TODO: Edit Request / Status / Less Errors / Email? / Notification? / Before PetSitter Response

// TODO: Cancel Request / Email? / Policy

// TODO: View Requests With Different Filters -> Accepted - Pending - Rejected - Canceled - ... / Different Sorts / Pagination

// TODO: Policy: Cancel - Price - calendar (warning)
