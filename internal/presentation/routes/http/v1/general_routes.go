package httpv1

import (
	"hona/backend/internal/application/service"
	"hona/backend/internal/infrastructure/database"
	"hona/backend/internal/infrastructure/repository/postgres"
	"hona/backend/internal/presentation/controllers/v1/general"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(v1 *gin.RouterGroup) {
	db := database.NewPostgresDatabase()
	r := postgres.NewUserRepository(db)
	s := service.NewGeneralService(r)
	gc := general.NewGeneralController(s)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", gc.Login)
	}
}
