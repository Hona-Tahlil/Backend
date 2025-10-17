package httpv1

import (
	"hona/backend/internal/application/service"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
	"hona/backend/internal/presentation/controllers/v1/general"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(v1 *gin.RouterGroup) {
	db := persistence.NewPostgresDatabase()
	r := postgres.NewUserRepository(db.GetDB())
	s := service.NewGeneralService(r)
	gc := general.NewGeneralController(s)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", gc.Login)
		auth.POST("/register", gc.Register)
		auth.POST("/verify", gc.VerifyEmail)
		auth.POST("/forgotpassword", gc.ForgotPassword)

	}
}
