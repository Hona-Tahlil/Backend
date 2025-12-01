package controllers

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/domain/exceptions"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

func GetTranslator(ctx *gin.Context, key string) ut.Translator {
	trans, _ := ctx.Get(key)
	return trans.(ut.Translator)
}

func SetRefreshTokenCookie(ctx *gin.Context, refreshToken string, expireTime int) {
	ctx.SetCookie(
		bootstrap.Run().Constants.Context.RefreshToken,
		refreshToken,
		expireTime,
		"/",
		"",
		true,
		true,
	)
}

func GetRefreshTokenCookie(ctx *gin.Context) (refreshToken string) {
	refreshToken, err := ctx.Cookie(bootstrap.Run().Constants.Context.RefreshToken)
	if err != nil {
		invalidTokenErr := exceptions.NewInvalidTokenError()
		panic(invalidTokenErr)
	}
	return
}

func GetID(ctx *gin.Context) uint {
	id, ok := ctx.Get(bootstrap.Run().Constants.Context.ID)
	if !ok {
		panic(exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User))
	}
	ID, _ := id.(uint)
	return ID
}
