package admin

import (
	"hona/backend/internal/application/usecase"
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
	controllers.Respond(ctx, 200, msg, res)
}
