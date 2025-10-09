package general

import (
	"hona/backend/internal/application/dto/login"
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralController struct {
	generalService *service.GeneralService
}

func NewGeneralController(generalService *service.GeneralService) *GeneralController {
	return &GeneralController{
		generalService: generalService,
	}
}

func (gc *GeneralController) Login(ctx *gin.Context) {
	type loginParams struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	params := controllers.Receive[loginParams](ctx)
	loginInfo := login.LoginRequest{
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
