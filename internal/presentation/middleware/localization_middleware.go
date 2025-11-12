package middleware

import (
	"net/http"

	"hona/backend/bootstrap"
	domaintranslation "hona/backend/internal/domain/translation"
	"hona/backend/internal/infrastructure/translation"

	"github.com/gin-gonic/gin"
)

type LocalizationMiddleware struct {
	translator domaintranslation.Translator
}

func NewLocalizationMiddleware() *LocalizationMiddleware {
	return &LocalizationMiddleware{
		translator: translation.NewTranslator(),
	}
}

func GetLocale(request *http.Request) string {
	return request.Header.Get(bootstrap.Run().Constants.Context.AcceptLanguage)
}

func (lm *LocalizationMiddleware) AddTranslator(ctx *gin.Context) {
	locale := GetLocale(ctx.Request)

	trans := lm.translator.GetTranslator(locale)

	ctx.Set(bootstrap.Run().Constants.Context.Translator, trans)

	ctx.Next()
}
