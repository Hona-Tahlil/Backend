package user

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type UserRequestController struct {
	requestService usecase.RequestService
}

func NewUserRequestController(requestService usecase.RequestService) *UserRequestController {
	return &UserRequestController{
		requestService: requestService,
	}

}

func (rc *UserRequestController) GetCreateRequestInfo(ctx *gin.Context) {
	type Params struct {
		PetSitterUserID uint `form:"petSitterUserID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)

	info := request.GetCreateRequestInfoRequest{
		UserID:          UserID,
		PetSitterUserID: params.PetSitterUserID,
	}
	res, err := rc.requestService.GetCreateRequestInfo(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (rc *UserRequestController) CreateRequest(ctx *gin.Context) {
	type CalendarSlot struct {
		Date  time.Time    `json:"date" binding:"required"`
		Slots []enums.Slot `json:"slots" binding:"required"`
	}
	type AddressInfo struct {
		ProvinceName  enums.Province `json:"provinceName" binding:"required"`
		CityName      enums.City     `json:"cityName" binding:"required"`
		StreetAddress string         `json:"streetAddress" binding:"required"`
		HouseNumber   uint           `json:"houseNumber" binding:"required"`
		Unit          uint           `json:"unit" binding:"required"`
		PostalCode    *string        `json:"postalCode"`
	}
	type Params struct {
		PetSitterUserID uint           `json:"petSitterUserID" validate:"required"`
		CalenderSlots   []CalendarSlot `json:"calendarSlots" validate:"required,min=1" binding:"dive"`
		PetIDs          []uint         `json:"petIDs" validate:"required"`
		Notes           *string        `json:"notes"`
		AddressInfo     *AddressInfo   `json:"addressInfo" validate:"omitempty"`
		AddressID       *uint          `json:"addressID"`
		ServiceID       uint           `json:"serviceID" validate:"required"`
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
		AddressInfo:     (*address.AddressInfo)(params.AddressInfo),
		AddressID:       params.AddressID,
		ServiceID:       params.ServiceID,
	}
	if err := rc.requestService.CreateRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (rc *UserRequestController) EditRequest(ctx *gin.Context) {
	type CalendarSlot struct {
		Date  time.Time    `json:"date" binding:"required"`
		Slots []enums.Slot `json:"slots" binding:"required"`
	}
	type AddressInfo struct {
		ProvinceName  enums.Province `json:"provinceName" binding:"required"`
		CityName      enums.City     `json:"cityName" binding:"required"`
		StreetAddress string         `json:"streetAddress" binding:"required"`
		HouseNumber   uint           `json:"houseNumber" binding:"required"`
		Unit          uint           `json:"unit" binding:"required"`
		PostalCode    *string        `json:"postalCode"`
	}
	type Params struct {
		RequestID     uint           `json:"requestID" validate:"required"`
		CalenderSlots []CalendarSlot `json:"calendarSlots" validate:"required,min=1" binding:"dive"`
		PetIDs        []uint         `json:"petIDs" validate:"required"`
		Notes         *string        `json:"notes"`
		AddressInfo   *AddressInfo   `json:"addressInfo" validate:"omitempty"`
		AddressID     *uint          `json:"addressID" validate:"required"`
		ServiceID     uint           `json:"serviceID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)
	slots := make([]request.RequestCalendarSlotRequest, len(params.CalenderSlots))
	for i, slot := range params.CalenderSlots {
		slots[i] = request.RequestCalendarSlotRequest{
			Date:  slot.Date,
			Slots: slot.Slots,
		}
	}
	info := request.EditRequestRequest{
		RequestID:     params.RequestID,
		UserID:        UserID,
		CalenderSlots: slots,
		PetIDs:        params.PetIDs,
		Notes:         params.Notes,
		AddressInfo:   (*address.AddressInfo)(params.AddressInfo),
		AddressID:     params.AddressID,
		ServiceID:     params.ServiceID,
	}

	if err := rc.requestService.EditRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (rc *UserRequestController) CancelRequest(ctx *gin.Context) {
	type Params struct {
		RequestID uint `json:"requestID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)

	info := request.CancelRequestRequest{
		RequestID: params.RequestID,
		UserID:    UserID,
	}

	if err := rc.requestService.CancelRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (rc *UserRequestController) GetRequestFullData(ctx *gin.Context) {
	type Params struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)

	info := request.GetRequestFullDataRequest{
		RequestID: params.RequestID,
	}
	res, err := rc.requestService.GetRequestFullData(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

// TODO: View Requests With Different Filters -> Accepted - Pending - Rejected - Canceled - ... / Different Sorts / Pagination
