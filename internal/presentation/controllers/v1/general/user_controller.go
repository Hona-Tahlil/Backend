package general

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/user"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	userService *service.UserService
}

func NewGeneralUserController(userService *service.UserService) *GeneralUserController {
	return &GeneralUserController{
		userService: userService,
	}

}

var successMessages = bootstrap.Run().Constants.SuccessMessages

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
	fmt.Println("✅ 1")

	res, refreshToken, expireTime, err := gc.userService.Login(loginInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("✅ 2")

	controllers.SetRefreshTokenCookie(ctx, refreshToken, expireTime)
	fmt.Println("✅ 3")

	msg := controllers.Message{
		Text: successMessages.Login,
	}
	fmt.Printf("LOGIN RES TYPE: %T\n", res)
	fmt.Printf("LOGIN RES VALUE: %#v\n", *res)
	
	controllers.Respond(ctx, 200, msg, *res)
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
	if err := gc.userService.Register(registerInfo); err != nil {
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
	if err := gc.userService.VerifyEmail(verifyOTPInfo); err != nil {
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
	if err := gc.userService.ForgotPassword(forgotPasswordInfo); err != nil {
		panic(err)
	}

	// trans := controller.GetTranslator(ctx, GeneralUserController.constants.Context.Translator)
	// message, _ := trans.Translate("successMessage.forgotPassword")
	// controller.Response(ctx, 200, message, nil)
}

func (gc *GeneralUserController) RefreshTokens(ctx *gin.Context) {
	RefreshToken := controllers.GetRefreshTokenCookie(ctx)

	refreshTokenInfo := rbac.RefreshTokenRequest{
		RefreshToken: RefreshToken,
	}

	res, refreshToken, expireTime, err := gc.userService.RefreshTokens(refreshTokenInfo)
	if err != nil {
		panic(err)
	}

	controllers.SetRefreshTokenCookie(ctx, refreshToken, expireTime)

	msg := controllers.Message{
		Text:   "successMessage.refreshToken",
		Params: []string{},
	}
	controllers.Respond(ctx, 200, msg, *res)
}
