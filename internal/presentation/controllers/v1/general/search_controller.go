package general

import (
	"hona/backend/internal/application/dto/general"
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

	offset, limit := controllers.GetOffsetLimit(params.Page, params.Count, 1, 10)

	// Convert params to DTO (no DSL conversion here)
	filters := make([]general.Filter, len(params.Filters))
	for i, f := range params.Filters {
		filters[i] = general.Filter{
			Field: f.Field,
			Op:    f.Op,
			Value: f.Value,
		}
	}

	sorts := make([]general.Sort, len(params.Sorts))
	for i, s := range params.Sorts {
		sorts[i] = general.Sort{
			Field: s.Field,
			Dir:   s.Dir,
		}
	}

	req := petsitter.SearchPetSittersRequest{
		Offset:  offset,
		Limit:   limit,
		Filters: filters,
		Sorts:   sorts,
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
