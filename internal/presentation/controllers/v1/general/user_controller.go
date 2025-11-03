package general

import (
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	generalService usecase.UserService
}

func NewGeneralUserController(generalService usecase.UserService) *GeneralUserController {
	return &GeneralUserController{
		generalService: generalService,
	}
}

func (gc *GeneralUserController) Login(ctx *gin.Context) {
	type loginParams struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	params := controllers.Receive[loginParams](ctx)
	loginInfo := user.LoginRequest{
		Email:    params.Email,
		Password: params.Password,
	}

	res := gc.generalService.Login(loginInfo)

	msg := controllers.Message{
		Text:   "success.login",
		Params: []string{},
	}
	controllers.Respond(ctx, 200, msg, res)
}

func (gc *GeneralUserController) Register(ctx *gin.Context) {
	type registerParams struct {
		FirstName string `json:"firstName" validate:"required"`
		LastName  string `json:"lastName" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required,min=8,max=64"`
	}
	params := controllers.Receive[registerParams](ctx)
	registerInfo := user.RegisterRequest{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}
	if err := gc.generalService.Register(registerInfo); err != nil {
		panic(err)
	}
	msg := controllers.Message{
		Text:   "success.userRegister",
		Params: []string{},
	}
	controllers.Respond(ctx, 200, msg, nil)
}

// func (gc *GeneralUserController) VerifyEmail(ctx *gin.Context) {
// 	type verifyPhoneParams struct {
// 		Email string `json:"phone" validate:"required,e164"`
// 		OTP   string `json:"otp" validate:"required"`
// 	}
// 	params := controllers.Receive[verifyPhoneParams](ctx)
// 	verifyOTPInfo := user.VerifyEmailRequest{
// 		Email: params.Email,
// 		OTP:   params.OTP,
// 	}
// 	if err := gc.generalService.VerifyEmail(verifyOTPInfo); err != nil {
// 		panic(err)
// 	}
// 	msg := controllers.Message{
// 		Text:   "success.emailVerification",
// 		Params: []string{},
// 	}
// 	controllers.Respond(ctx, 200, msg, nil)

// }

// func (gc *GeneralUserController) ForgotPassword(ctx *gin.Context) {
// 	type forgotPasswordParams struct {
// 		Email string `json:"Email" validate:"required,e164"`
// 	}
// 	params := controllers.Receive[forgotPasswordParams](ctx)
// 	forgotPasswordInfo := user.ForgotPasswordRequest{
// 		Email: params.Email,
// 	}
// 	if err := gc.generalService.ForgotPassword(forgotPasswordInfo); err != nil {
// 		panic(err)
// 	}

// }
