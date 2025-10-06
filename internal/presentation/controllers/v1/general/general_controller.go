package general

import "github.com/gin-gonic/gin"

func Login(ctx *gin.Context) {
	type loginParams struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// TODO: Receive

	// TODO: Call Service

	// TODO: Respond
}
