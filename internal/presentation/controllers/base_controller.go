package controllers

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

func GetTranslator(ctx *gin.Context, key string) ut.Translator {
	trans, _ := ctx.Get(key)
	return trans.(ut.Translator)
}
