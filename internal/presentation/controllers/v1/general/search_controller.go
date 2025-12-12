package general

import (
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/service"
	"hona/backend/internal/infrastructure/dsl"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralPetSitterController struct {
	petSitterService *service.PetSitterService
}

func NewGeneralPetSitterController(petSitterService *service.PetSitterService) *GeneralPetSitterController {
	return &GeneralPetSitterController{
		petSitterService: petSitterService,
	}
}

// @Summary Search pet sitters
// @Description Search for pet sitters by name, bio, or location
// @Tags Pet Sitters
// @Param q query string true "Search query"
// @Param province query string false "Filter by province"
// @Param page query int false "Page number" default(1)
// @Param count query int false "Results per page" default(10)
// @Param sort_by query string false "Sort by: rating or experience" default(id)
// @Success 200 {object} petsitter.SearchPetSittersResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /v1/pet-sitters/search [get]
// type Filter struct {
// 	Field string
// 	Op    string // =, !=, >, <, >=, <=, LIKE, IN
// 	Value any
// }

// type Filters struct {
// 	Items []Filter
// }

// func (gc *GeneralPetSitterController) SearchPetSitters(c *gin.Context) {
//     query := c.Query("q")
//     if query == "" {
//         response.Error(c, http.StatusBadRequest, "Search query is required")
//         return
//     }

//     page := 1
//     if p := c.Query("page"); p != "" {
//         if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
//             page = pageNum
//         }
//     }

//     count := 10
//     if cnt := c.Query("count"); cnt != "" {
//         if countNum, err := strconv.Atoi(cnt); err == nil && countNum > 0 && countNum <= 100 {
//             count = countNum
//         }
//     }

//     req := petsitter.SearchPetSittersRequest{
//         Query:    query,
//         Province: c.Query("province"),
//         Page:     page,
//         Count:    count,
//         SortBy:   c.Query("sort_by"),
//     }

//     result, err := gc.petSitterService.SearchPetSitters(req)
//     if err != nil {
//         response.Error(c, http.StatusInternalServerError, "Failed to search pet sitters")
//         return
//     }

//     response.Success(c, http.StatusOK, "Pet sitters found", result)
// }

func (gc *GeneralPetSitterController) SearchPetSitters(ctx *gin.Context) {
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
	

	filters := make([]dsl.Filter,0)
	for _, f := range params.Filters {
		filter  := dsl.Filter {
			Field: f.Field,
			Op: f.Op,
			Value: f.Value,
		}
		filters = append(filters, filter)
	}

	sorts := make([]dsl.Sort,0)
	for _,s := range params.Sorts {
		sort := dsl.Sort {
			Field: s.Field,
			Dir: s.Dir,
		}
		sorts = append(sorts, sort)
	}

	req := petsitter.SearchPetSittersRequest{
		Offset:    offset,
		Limit:   limit,
		Filters: filters,
		Sorts:   sorts,
	}

	petsitters ,count , err := gc.petSitterService.SearchPetSitters(req)
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
