package routes

import (
	httpv1 "hona/backend/internal/presentation/routes/http/v1"
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRoutes(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.LoggingMiddleware.LogRequests)
	ginEngine.Use(app.Middlewares.CORSMiddleware.CORS())
	ginEngine.Use(app.Middlewares.LocalizationMiddleware.AddTranslator)
	ginEngine.Use(app.Middlewares.RecoveryMiddleware.Recover)
	ginEngine.Use(app.Middlewares.Prometheus.PrometheusMiddleware)

	ginEngine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := ginEngine.Group("/v1")
	httpv1.SetUpGeneralRoutes(v1, app)
	httpv1.SetUpAdminRoutes(v1, app)
	httpv1.SetUpUserRoutes(v1, app)
	httpv1.SetUpPetSitterRoutes(v1, app)
}
