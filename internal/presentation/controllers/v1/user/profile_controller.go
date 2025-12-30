package user

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/presentation/controllers"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type UserProfileController struct {
	userService usecase.UserService
}

func NewUserProfileController(userService usecase.UserService) *UserProfileController {
	return &UserProfileController{
		userService: userService,
	}
}

func (pc *UserProfileController) GetProfile(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	info := user.GetProfileRequest{
		UserID: UserID,
	}
	res, err := pc.userService.GetProfile(info)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (pc *UserProfileController) GetIdentity(ctx *gin.Context) {
	UserID := controllers.GetID(ctx)
	info := user.GetIdentityRequest{
		UserID: UserID,
	}
	res, err := pc.userService.GetIdentity(info)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, *res)
}

func (pc *UserProfileController) UpdateProfile(ctx *gin.Context) {
	type UpdateProfileParams struct {
		FirstName     string         `form:"firstName" validate:"required"`
		LastName      string         `form:"lastName" validate:"required"`
		Phone         *string        `form:"phone"`
		Gender        enums.Gender   `form:"gender" validate:"required"`
		BirthDate     *time.Time     `form:"birthDate"`
		Province      enums.Province `form:"province"`
		City          enums.City     `form:"city"`
		StreetAddress string         `form:"streetAddress"`
		HouseNumber   uint           `form:"houseNumber"`
		Unit          uint           `form:"unit"`
		PostalCode    *string        `form:"postalCode"`
	}

	file, err := ctx.FormFile(bootstrap.Run().Env.Storage.Buckets.UserProfilePic)
	if err != nil {
		// ve := exceptions.NewValidationErrors()
		// ve.AddError(bootstrap.Run().Env.Storage.Buckets.UserProfilePic, bootstrap.Run().Constants.ErrorTags.Required)
		// panic(ve)
		log.Println("no file")
	}

	params := controllers.Receive[UpdateProfileParams](ctx)
	UserID := controllers.GetID(ctx)
	info := user.UpdateProfileRequest{
		UserID:        UserID,
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
		ProfilePic:    file,
	}

	if err := pc.userService.UpdateProfile(info); err != nil {
		panic(err)
	}

	msg := controllers.Message{
		Text: successMessages.UpdateProfile,
	}
	controllers.Respond(ctx, 200, msg, nil)
}
