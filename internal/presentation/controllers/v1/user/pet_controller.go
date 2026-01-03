package user

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type UserPetController struct {
	petService usecase.PetService
}

func NewUserPetController(petService usecase.PetService) *UserPetController {
	return &UserPetController{
		petService: petService,
	}
}

var successMessages = bootstrap.Run().Constants.SuccessMessages

func (uc *UserPetController) AddPet(ctx *gin.Context) {
	type AddPetParams struct {
		Name      string     `form:"name" validate:"required,min=1,max=100"`
		Kind      uint       `form:"kind" validate:"required,min=1,max=18"`
		Species   uint       `form:"species" validate:"required,min=1,max=67"`
		BirthDate *time.Time `form:"birthDate"`
		IsAdult   bool       `form:"isAdult" validate:"omitempty"`
		Gender    uint       `form:"gender" validate:"required,min=1,max=3"`
		Weight    *float32   `form:"weight" validate:"omitempty,min=0.1,max=500"`
		AboutPet  *string    `form:"aboutPet" validate:"omitempty,max=10000"`
	}
	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.PetProfilePic)
	if err != nil {
		file = nil
	}
	params := controllers.Receive[AddPetParams](ctx)
	UserID := controllers.GetID(ctx)
	AddPetInfo := pet.AddPetRequest{
		UserID:     UserID,
		Name:       params.Name,
		Kind:       enums.PetKind(params.Kind),
		Species:    enums.Species(params.Species),
		BirthDate:  params.BirthDate,
		IsAdult:    params.IsAdult,
		Gender:     enums.PetGender(params.Gender),
		Weight:     params.Weight,
		AboutPet:   params.AboutPet,
		ProfilePic: file,
	}
	err = uc.petService.AddPet(AddPetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: successMessages.AddPet,
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) UpdatePet(ctx *gin.Context) {
	type UpdatePetParams struct {
		ID        uint       `form:"id" validate:"required"`
		Name      string     `form:"name" validate:"required,min=1,max=100"`
		Kind      uint       `form:"kind" validate:"required,min=1,max=18"`
		Species   uint       `form:"species" validate:"required,min=1,max=67"`
		BirthDate *time.Time `form:"birthDate" validate:"omitempty"`
		IsAdult   bool       `form:"isAdult" validate:"omitempty"`
		Gender    uint       `form:"gender" validate:"required,min=1,max=3"`
		Weight    *float32   `form:"weight" validate:"omitempty,min=0.1,max=500"`
		AboutPet  *string    `form:"aboutPet" validate:"omitempty,max=10000"`
	}
	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.PetProfilePic)
	if err != nil {
		file = nil
	}
	params := controllers.Receive[UpdatePetParams](ctx)
	UserID := controllers.GetID(ctx)
	UpdatePetInfo := pet.UpdatePetRequest{
		UserID:     UserID,
		ID:         params.ID,
		Name:       params.Name,
		Kind:       enums.PetKind(params.Kind),
		Species:    enums.Species(params.Species),
		BirthDate:  params.BirthDate,
		IsAdult:    params.IsAdult,
		Gender:     enums.PetGender(params.Gender),
		Weight:     params.Weight,
		AboutPet:   params.AboutPet,
		ProfilePic: file,
	}

	err = uc.petService.UpdatePet(UpdatePetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: successMessages.UpdatePet,
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) RemovePet(ctx *gin.Context) {
	type RemovePetParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[RemovePetParams](ctx)
	UserID := controllers.GetID(ctx)
	RemovePetInfo := pet.RemovePetRequest{
		UserID: UserID,
		ID:     params.ID,
	}
	err := uc.petService.RemovePet(RemovePetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: successMessages.RemovePet,
	}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) GetPetsBasicData(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	GetPetsBasicDataInfo := pet.GetPetsBasicDataRequest{
		UserID: UserID,
	}
	res, err := uc.petService.GetPetsBasicData(GetPetsBasicDataInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (us *UserPetController) GetPetFullData(ctx *gin.Context) {
	type GetPetFullDataParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[GetPetFullDataParams](ctx)
	GetPetFullDataInfo := pet.GetPetFullDataRequest{
		ID: params.ID,
	}
	res, err := us.petService.GetPetFullData(GetPetFullDataInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}
