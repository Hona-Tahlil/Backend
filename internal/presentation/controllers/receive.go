package controllers

import "github.com/gin-gonic/gin"

func Receive[T any](ctx *gin.Context) T {
	var params T

	if err := ctx.ShouldBind(&params); err != nil {
		// TODO: Proper Error Handling
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		// TODO: Proper Error Handling
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		// TODO: Proper Error Handling
	}

	// TODO: Validation

	return params
}
