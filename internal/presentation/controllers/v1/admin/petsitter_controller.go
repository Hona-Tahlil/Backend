package admin

import (
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type AdminPetSitterController struct {
	petSitterService usecase.PetSitterService
}

func NewAdminPetSitterController(petSitterService usecase.PetSitterService) *AdminPetSitterController {
	return &AdminPetSitterController{
		petSitterService: petSitterService,
	}
}

func (apc *AdminPetSitterController) ListAllPetSitters(ctx *gin.Context) {
	type ListPetSittersParams struct {
		Page  int `form:"page" validate:"min=1"`
		Count int `form:"count" validate:"min=1,max=100"`
	}
	params := controllers.Receive[ListPetSittersParams](ctx)

	// Set defaults
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Count == 0 {
		params.Count = 10
	}

	res, err := apc.petSitterService.GetAllPetSitters(params.Page, params.Count)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (apc *AdminPetSitterController) SearchPetSitters(ctx *gin.Context) {
	type Filter struct {
		Field string `json:"field" validate:"required"`
		Op    string `json:"op"    validate:"required,oneof== != > < >= <= LIKE IN"`
		Value any    `json:"value" validate:"required"`
	}
	type Sort struct {
		Field string `json:"field" validate:"required"`
		Dir   string `json:"dir"   validate:"required,oneof=ASC DESC"`
	}
	type SearchPetSittersParams struct {
		Page    int      `json:"page" validate:"omitempty,min=1"`
		Count   int      `json:"count" validate:"omitempty,min=1,max=100"`
		Filters []Filter `json:"filters" validate:"omitempty,dive"`
		Sorts   []Sort   `json:"sorts" validate:"omitempty,dive"`
	}
	params := controllers.Receive[SearchPetSittersParams](ctx)
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

	req := petsitter.AdminSearchPetSittersRequest{
		Offset:  offset,
		Limit:   limit,
		Filters: controllers.ToFilters(filterParams),
		Sorts:   controllers.ToSorts(sortParams),
	}
	items, total, err := apc.petSitterService.SearchPetSittersForAdmin(req)
	if err != nil {
		panic(err)
	}

	data := controllers.NewPaginatedResponse(items, total, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}

func (apc *AdminPetSitterController) GetPetSitterDetails(ctx *gin.Context) {
	type Params struct {
		PetSitterUserID uint `uri:"petsitterUserID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	info := petsitter.GetPetSitterDetailsRequest{
		PetSitterUserID: params.PetSitterUserID,
	}

	res, err := apc.petSitterService.GetPetSitterDetails(info)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (apc *AdminPetSitterController) ChangePetSitterStatus(ctx *gin.Context) {
	type Params struct {
		PetSitterUserID uint                  `json:"petsitterUserID" validate:"required"`
		Status          enums.PetSitterStatus `json:"status" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)
	info := petsitter.ChangePetSitterStatusRequest{
		PetSitterUserID: params.PetSitterUserID,
		Status:          params.Status,
	}

	err := apc.petSitterService.ChangePetSitterStatus(info)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}
