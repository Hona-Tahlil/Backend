package general

import "github.com/gin-gonic/gin"

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

	// TODO: Receive

	// TODO: Call Service

	// TODO: Respond
}
