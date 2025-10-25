package routes

import (
	httpv1 "hona/backend/internal/presentation/routes/http/v1"
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.LocalizationMiddleware.AddTranslator)
	ginEngine.Use(app.Middlewares.RecoveryMiddleware.Recover)

	v1 := ginEngine.Group("/v1")
	httpv1.SetUpGeneralRoutes(v1, app)
}
