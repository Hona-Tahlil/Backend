package petsitter

import (
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type PetSitterRequestController struct {
	requestService usecase.RequestService
}

func NewPetSitterRequestController(requestService usecase.RequestService) *PetSitterRequestController {
	return &PetSitterRequestController{
		requestService: requestService,
	}
}

func (pc *PetSitterRequestController) GetRequestFullData(ctx *gin.Context) {
	type Params struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)

	info := request.GetRequestFullDataRequest{
		RequestID: params.RequestID,
	}
	res, err := pc.requestService.GetRequestFullData(info)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

// TODO: View Requests With Different Filters -> Accepted - Pending - Rejected - Canceled - ... / Different Sorts / Pagination

func (pc *PetSitterRequestController) CancelRequest(ctx *gin.Context) {
	type Params struct {
		RequestID uint `json:"requestID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)

	info := request.CancelRequestRequest{
		RequestID: params.RequestID,
		UserID:    UserID,
	}

	if err := pc.requestService.CancelRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterRequestController) RespondToRequest(ctx *gin.Context) {
	type Params struct {
		GetTime   time.Time `json:"getTime" validate:"required"`
		Accept    bool      `json:"accept" validate:"required"`
		RequestID uint      `json:"requestID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	UserID := controllers.GetID(ctx)
	info := request.RespondToRequestRequest{
		RequestID: params.RequestID,
		Accept:    params.Accept,
		UserID:    UserID,
	}

	if err := pc.requestService.RespondToRequest(info); err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}
