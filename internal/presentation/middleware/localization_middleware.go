package middleware

import (
	"net/http"

	"hona/backend/internal/infrastructure/translation"

	"github.com/gin-gonic/gin"
)

// ? what should it have?
type LocalizationMiddleware struct {
}

func NewLocalizationMiddleware() *LocalizationMiddleware {
	return &LocalizationMiddleware{}
}

func GetLocale(request *http.Request) string {
	return request.Header.Get("Accept-Language")
}

func (lm *LocalizationMiddleware) AddTranslator(ctx *gin.Context) {
	locale := GetLocale(ctx.Request)

	trans := translation.GetTranslator(locale)

	ctx.Set("translator", trans)

	ctx.Next()
}
