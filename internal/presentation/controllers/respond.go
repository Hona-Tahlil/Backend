package controllers

import "github.com/gin-gonic/gin"

type Message struct {
	Text   string
	Params []string
}

func Respond[T Message | []Message](ctx *gin.Context, statusCode int, messages T, data interface{}) {
	// TODO: Get Translator

	// TODO: Implement
}
