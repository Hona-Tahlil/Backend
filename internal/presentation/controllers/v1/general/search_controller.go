package general

import (
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralSearchController struct {
	petSitterService *service.PetSitterService
}

func NewGeneralSearchController(petSitterService *service.PetSitterService) *GeneralSearchController {
	return &GeneralSearchController{
		petSitterService: petSitterService,
	}
}

func (gc *GeneralSearchController) SearchPetSitters(ctx *gin.Context) {
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


	req := petsitter.SearchPetSittersRequest{
		Offset:  offset,
		Limit:   limit,
		Filters: controllers.ToFilters(filterParams),
		Sorts:   controllers.ToSorts(sortParams),
	}

	petsitters, count, err := gc.petSitterService.SearchPetSitters(req)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(petsitters, count, offset, limit)

	msg := controllers.Message{
		Text:   "success.petSitterSearch",
		Params: []string{},
	}
	controllers.Respond(ctx, 200, msg, data)
}

func (gc *GeneralSearchController) GetPetSitterProfile(ctx *gin.Context) {
	type Params struct {
		PetSitterID uint `uri:"petSitterID" validate:"required"`
	}
	params := controllers.Receive[Params](ctx)

	req := petsitter.GetPetSitterProfileRequest{
		PetSitterID: params.PetSitterID,
	}
	res, err := gc.petSitterService.GetPetSitterProfile(req)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
