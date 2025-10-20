package middleware

import (
	"net/http"

	"hona/backend/bootstrap"
	"hona/backend/internal/infrastructure/translation"

	"github.com/gin-gonic/gin"
)

type LocalizationMiddleware struct {
}

func NewLocalizationMiddleware() *LocalizationMiddleware {
	return &LocalizationMiddleware{}
}

func GetLocale(request *http.Request) string {
	return request.Header.Get(bootstrap.Run().Constants.Context.AcceptLanguage)
}

func (lm *LocalizationMiddleware) AddTranslator(ctx *gin.Context) {
	locale := GetLocale(ctx.Request)

	trans := translation.GetTranslator(locale)

	ctx.Set(bootstrap.Run().Constants.Context.Translator, trans)

	ctx.Next()
}
