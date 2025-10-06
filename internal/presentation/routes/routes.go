package routes

import (
	"hona/backend/internal/presentation/middleware"
	httpv1 "hona/backend/internal/presentation/routes/http/v1"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(ginEngine *gin.Engine) {
	lm := middleware.NewLocalizationMiddleware()
	rm := middleware.NewRecoveryMiddleware()

	ginEngine.Use(lm.AddTranslator)
	ginEngine.Use(rm.Recover)

	v1 := ginEngine.Group("/v1")
	httpv1.SetUpRoutes(v1)
}
