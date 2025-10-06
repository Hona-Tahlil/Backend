package routes

import (
	httpv1 "hona/backend/internal/presentation/routes/http/v1"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(ginEngine *gin.Engine) {
	v1 := ginEngine.Group("/v1")
	httpv1.SetUpRoutes(v1)
}
