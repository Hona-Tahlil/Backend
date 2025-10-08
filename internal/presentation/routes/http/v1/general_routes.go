package httpv1

import (
	"hona/backend/internal/application/service"
	"hona/backend/internal/presentation/controllers/v1/general"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(v1 *gin.RouterGroup) {
	s := service.NewGeneralService()
	gc := general.NewGeneralController(s)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", gc.Login)
	}
}
