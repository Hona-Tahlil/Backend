package general

import (
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralPetController struct {
	petService usecase.PetService
}

func NewGeneralPetController(petService usecase.PetService) *GeneralPetController {
	return &GeneralPetController{
		petService: petService,
	}
}

func (pc *GeneralPetController) GetAllPetKinds(ctx *gin.Context) {
	res := pc.petService.GetAllPetKinds()

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *GeneralPetController) GetPetKindSpecies(ctx *gin.Context) {
	type Params struct {
		PetKind enums.PetKind `uri:"petKind" validate:"required,min=1,max=18"`
	}
	params := controllers.Receive[Params](ctx)
	info := pet.GetPetKindSpecies{
		Num: params.PetKind,
	}
	res := pc.petService.GetPetKindSpecies(info)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
