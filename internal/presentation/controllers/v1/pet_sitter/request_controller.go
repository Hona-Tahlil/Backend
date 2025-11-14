package petsitter

import (
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterRequestController struct {
	requestService *service.RequestService
}

func NewPetSitterRequestController(requestService *service.RequestService) *PetSitterRequestController {
	return &PetSitterRequestController{
		requestService: requestService,
	}
}

func (pc *PetSitterRequestController) GetRequestFullData(ctx *gin.Context) {
	type Params struct {
		RequestID uint `uri:"requestID"`
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

// TODO: Email? / Policy
func (pc *PetSitterRequestController) CancelRequest(ctx *gin.Context) {
	type Params struct {
		RequestID uint `json:"requestID"`
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
		Accept    bool `json:"accept"`
		RequestID uint `json:"requestID"`
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
