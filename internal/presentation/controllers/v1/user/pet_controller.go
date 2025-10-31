package user

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/service"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type UserPetController struct {
	petService *service.PetService
}

func NewUserPetController(petService *service.PetService) *UserPetController {
	return &UserPetController{
		petService: petService,
	}
}

func (uc *UserPetController) AddPet(ctx *gin.Context) {
	type AddPetParams struct {
		Name      string          `json:"name" validate:"required,min=1,max=100"`
		Kind      enums.PetKind   `json:"kind" validate:"required,min=1,max=20"`     // TODO: exact number for max
		Species   enums.Species   `json:"species" validate:"required,min=1,max=100"` // TODO: exact number for max
		BirthDate *time.Time      `json:"birthDate" validate:"omitempty,datetime"`
		IsAdult   bool            `json:"isAdult" validate:"omitempty"`
		Gender    enums.PetGender `json:"gender" validate:"omitempty,min=1,max=3"`
		Weight    *float32        `json:"weight" validate:"omitempty,min=0.1,max=500"`
		AboutPet  *string         `json:"aboutPet" validate:"omitempty,max=10000"`
	}
	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.PetProfilePic)
	if err != nil {
		file = nil
	}
	params := controllers.Receive[AddPetParams](ctx)
	userID, _ := ctx.Get(bootstrap.Run().Constants.Context.ID)
	UserID, _ := userID.(uint)
	AddPetInfo := pet.AddPetRequest{
		UserID:     UserID,
		Name:       params.Name,
		Kind:       params.Kind,
		Species:    params.Species,
		BirthDate:  params.BirthDate,
		IsAdult:    params.IsAdult,
		Gender:     params.Gender,
		Weight:     params.Weight,
		AboutPet:   params.AboutPet,
		ProfilePic: file,
	}
	err = uc.petService.AddPet(AddPetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) UpdatePet(ctx *gin.Context) {
	type UpdatePetParams struct {
		ID        uint            `json:"id" validate:"required"`
		Name      string          `json:"name" validate:"required,min=1,max=100"`
		Kind      enums.PetKind   `json:"kind" validate:"required,min=1,max=20"`     // TODO: exact number for max
		Species   enums.Species   `json:"species" validate:"required,min=1,max=100"` // TODO: exact number for max
		BirthDate *time.Time      `json:"birthDate" validate:"omitempty,datetime"`
		IsAdult   bool            `json:"isAdult" validate:"omitempty"`
		Gender    enums.PetGender `json:"gender" validate:"omitempty,min=1,max=3"`
		Weight    *float32        `json:"weight" validate:"omitempty,min=0.1,max=500"`
		AboutPet  *string         `json:"aboutPet" validate:"omitempty,max=10000"`
	}
	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.PetProfilePic)
	if err != nil {
		file = nil
	}
	params := controllers.Receive[UpdatePetParams](ctx)
	UpdatePetInfo := pet.UpdatePetRequest{
		ID:         params.ID,
		Name:       params.Name,
		Kind:       params.Kind,
		Species:    params.Species,
		BirthDate:  params.BirthDate,
		IsAdult:    params.IsAdult,
		Gender:     params.Gender,
		Weight:     params.Weight,
		AboutPet:   params.AboutPet,
		ProfilePic: file,
	}

	err = uc.petService.UpdatePet(UpdatePetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) RemovePet(ctx *gin.Context) {
	type RemovePetParams struct {
		ID uint `uri:"id"`
	}
	params := controllers.Receive[RemovePetParams](ctx)
	RemovePetInfo := pet.RemovePetRequest{
		ID: params.ID,
	}
	err := uc.petService.RemovePet(RemovePetInfo)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (uc *UserPetController) GetPetsBasicData(ctx *gin.Context) {
	userID, _ := ctx.Get(bootstrap.Run().Constants.Context.ID)
	UserID, _ := userID.(uint)
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
