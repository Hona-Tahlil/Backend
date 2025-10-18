package general

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	generalService *service.UserService
	constants      *bootstrap.Constants
}

func NewGeneralUserController(generalService *service.UserService) *GeneralUserController {
	return &GeneralUserController{
		generalService: generalService,
	}
}

func (gc *GeneralUserController) Login(ctx *gin.Context) {
	type loginParams struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=8,max=64"`
		RememberMe bool   `json:"rememberMe"`
	}

	params := controllers.Receive[loginParams](ctx)
	loginInfo := user.LoginRequest{
		Email:      params.Email,
		Password:   params.Password,
		RememberMe: params.RememberMe,
	}

	res, refreshToken, err := gc.generalService.Login(loginInfo)
	if err != nil {
		panic(err)
	}

	ctx.SetCookie(
		"refreshToken",
		refreshToken,
		int(time.Hour.Seconds()*7*24),
		"/",
		"",
		true,
		true,
	)

	msg := controllers.Message{
		Text:   "successMessage.login",
		Params: []string{},
	}
	controllers.Respond(ctx, 200, msg, res)
}

func (gc *GeneralUserController) Register(ctx *gin.Context) {
	type registerParams struct {
		FirstName       string `json:"firstName" validate:"required"`
		LastName        string `json:"lastName" validate:"required"`
		Email           string `json:"email" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
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

}

func (gc *GeneralUserController) VerifyEmail(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Email string `json:"phone" validate:"required,e164"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controllers.Receive[verifyPhoneParams](ctx)
	verifyOTPInfo := user.VerifyEmailRequest{
		Email: params.Email,
		OTP:   params.OTP,
	}
	if err := gc.generalService.VerifyEmail(verifyOTPInfo); err != nil {
		panic(err)
	}

	// trans := controller.GetTranslator(ctx, GeneralUserController.constants.Context.Translator)
	// message, _ := trans.Translate("successMessage.phoneVerification")
	// controller.Response(ctx, 200, message, nil)
}

func (gc *GeneralUserController) ForgotPassword(ctx *gin.Context) {
	type forgotPasswordParams struct {
		Email string `json:"Email" validate:"required,e164"`
	}
	params := controllers.Receive[forgotPasswordParams](ctx)
	forgotPasswordInfo := user.ForgotPasswordRequest{
		Email: params.Email,
	}
	if err := gc.generalService.ForgotPassword(forgotPasswordInfo); err != nil {
		panic(err)
	}

	// trans := controller.GetTranslator(ctx, GeneralUserController.constants.Context.Translator)
	// message, _ := trans.Translate("successMessage.forgotPassword")
	// controller.Response(ctx, 200, message, nil)
}
