package petsitter

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterSkillsController struct {
	petSitterService usecase.PetSitterService
}

func NewPetSitterSkillsController(petSitterService usecase.PetSitterService) *PetSitterSkillsController {
	return &PetSitterSkillsController{
		petSitterService: petSitterService,
	}
}

func (pc *PetSitterSkillsController) GetPetKinds(ctx *gin.Context) {
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetPetKinds(petsitter.GetPetKindsRequest{UserID: userID})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterSkillsController) UpdatePetKinds(ctx *gin.Context) {
	type UpdatePetKindsParams struct {
		PetKinds []enums.PetKind `json:"pet_kinds" validate:"omitempty,dive,min=1"`
	}
	params := controllers.Receive[UpdatePetKindsParams](ctx)
	userID := controllers.GetID(ctx)
	err := pc.petSitterService.UpdatePetKinds(petsitter.UpdatePetKindsRequest{
		UserID:   userID,
		PetKinds: params.PetKinds,
	})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, nil)
}

func (pc *PetSitterSkillsController) GetServices(ctx *gin.Context) {
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetServices(petsitter.GetServicesRequest{UserID: userID})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterSkillsController) CreateService(ctx *gin.Context) {
	type CreateServiceParams struct {
		Type        enums.ServiceType `json:"type" validate:"required,min=1"`
		Price       uint              `json:"price" validate:"min=0"`
		Description *string           `json:"description" validate:"omitempty,max=10000"`
	}
	params := controllers.Receive[CreateServiceParams](ctx)
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.CreateService(petsitter.CreateServiceRequest{
		UserID:      userID,
		Type:        params.Type,
		Price:       params.Price,
		Description: params.Description,
	})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterSkillsController) UpdateService(ctx *gin.Context) {
	type UpdateServiceParams struct {
		ID          uint              `json:"id" validate:"required"`
		Type        enums.ServiceType `json:"type" validate:"required,min=1"`
		Price       uint              `json:"price" validate:"min=0"`
		Description *string           `json:"description" validate:"omitempty,max=10000"`
	}
	params := controllers.Receive[UpdateServiceParams](ctx)
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.UpdateService(petsitter.UpdateServiceRequest{
		UserID:      userID,
		ID:          params.ID,
		Type:        params.Type,
		Price:       params.Price,
		Description: params.Description,
	})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, res)
}

func (pc *PetSitterSkillsController) DeleteService(ctx *gin.Context) {
	type DeleteServiceParams struct {
		ID uint `uri:"serviceID" validate:"required"`
	}
	params := controllers.Receive[DeleteServiceParams](ctx)
	userID := controllers.GetID(ctx)
	err := pc.petSitterService.DeleteService(petsitter.DeleteServiceRequest{
		UserID: userID,
		ID:     params.ID,
	})
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, nil)
}
