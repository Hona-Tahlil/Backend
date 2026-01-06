package petsitter

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type PetSitterProfileController struct {
	petSitterService usecase.PetSitterService
}

func NewPetSitterProfileController(petSitterService usecase.PetSitterService) *PetSitterProfileController {
	return &PetSitterProfileController{
		petSitterService: petSitterService,
	}
}

func (pc *PetSitterProfileController) GetProfile(ctx *gin.Context) {
	userID := controllers.GetID(ctx)
	res, err := pc.petSitterService.GetPetSitterSelfProfile(petsitter.GetPetSitterSelfProfileRequest{UserID: userID})
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (pc *PetSitterProfileController) UpdateProfile(ctx *gin.Context) {
	type UpdateProfileParams struct {
		FirstName     string          `form:"firstName" validate:"required"`
		LastName      string          `form:"lastName" validate:"required"`
		Phone         *string         `form:"phone"`
		Gender        enums.Gender    `form:"gender" validate:"required"`
		BirthDate     *time.Time      `form:"birthDate"`
		Province      enums.Province  `form:"province"`
		City          enums.City      `form:"city"`
		StreetAddress string          `form:"streetAddress"`
		HouseNumber   uint            `form:"houseNumber"`
		Unit          uint            `form:"unit"`
		PostalCode    *string         `form:"postalCode"`
		Bio           *string         `form:"bio" validate:"omitempty"`
	}

	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.UserProfilePic)
	if err != nil {
		file = nil
	}

	params := controllers.Receive[UpdateProfileParams](ctx)
	userID := controllers.GetID(ctx)

	info := petsitter.UpdatePetSitterProfileRequest{
		UserID:        userID,
		FirstName:     params.FirstName,
		LastName:      params.LastName,
		Phone:         params.Phone,
		Gender:        params.Gender,
		BirthDate:     params.BirthDate,
		Province:      params.Province,
		City:          params.City,
		StreetAddress: params.StreetAddress,
		HouseNumber:   params.HouseNumber,
		Unit:          params.Unit,
		PostalCode:    params.PostalCode,
		Bio:           params.Bio,
		ProfilePic:    file,
	}

	if err := pc.petSitterService.UpdatePetSitterProfile(info); err != nil {
		panic(err)
	}

	msg := controllers.Message{Text: bootstrap.Run().Constants.SuccessMessages.Generic}
	controllers.Respond(ctx, 200, msg, nil)
}
