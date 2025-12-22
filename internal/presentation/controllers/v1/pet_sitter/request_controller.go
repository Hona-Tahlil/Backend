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

func (pc *PetSitterRequestController) SearchRequests(ctx *gin.Context) {
	type Filter struct {
		Field string `json:"field" validate:"required"`
		Op    string `json:"op"    validate:"required,oneof== != > < >= <= LIKE IN"`
		Value any    `json:"value" validate:"required"`
	}
	type Sort struct {
		Field string `json:"field" validate:"required"`
		Dir   string `json:"dir"   validate:"required,oneof=ASC DESC"`
	}
	type SearchRequestsParams struct {
		Page    int      `json:"page" validate:"omitempty,min=1"`
		Count   int      `json:"count" validate:"omitempty,min=1,max=100"`
		Filters []Filter `json:"filters" validate:"omitempty,dive"`
		Sorts   []Sort   `json:"sorts" validate:"omitempty,dive"`
	}

	params := controllers.Receive[SearchRequestsParams](ctx)
	userID := controllers.GetID(ctx)
	offset, limit := controllers.GetOffsetLimit(params.Page, params.Count)

	filterParams := make([]controllers.FilterParams, len(params.Filters))
	for i, f := range params.Filters {
		filterParams[i] = controllers.FilterParams{
			Field: f.Field,
			Op:    f.Op,
			Value: f.Value,
		}
	}

	sortParams := make([]controllers.SortParams, len(params.Sorts))
	for i, s := range params.Sorts {
		sortParams[i] = controllers.SortParams{
			Field: s.Field,
			Dir:   s.Dir,
		}
	}

	req := request.SearchPetSitterRequestsRequest{
		PetSitterUserID: userID,
		Offset:          offset,
		Limit:           limit,
		Filters:         controllers.ToFilters(filterParams),
		Sorts:           controllers.ToSorts(sortParams),
	}
	requests, count, err := pc.requestService.SearchPetSitterRequests(req)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(requests, count, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
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
