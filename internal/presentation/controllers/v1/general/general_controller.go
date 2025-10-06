package general

import (
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralController struct {
}

func NewGeneralController() *GeneralController {
	return &GeneralController{}
}

func (gc *GeneralController) Login(ctx *gin.Context) {
	type loginParams struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	controllers.Receive[loginParams](ctx)

	// TODO: Call Service

	// TODO: Respond
}
