package httpv1

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/service"
	"hona/backend/internal/infrastructure/mail"
	"hona/backend/internal/infrastructure/persistence"
	"hona/backend/internal/infrastructure/persistence/repository/redis"
	"hona/backend/internal/presentation/controllers/v1/general"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(v1 *gin.RouterGroup) {
	db := persistence.NewPostgresDatabase()
	uw := persistence.NewUnitOfWork(db.DB)
	rdb := persistence.NewRedisDatabase(&bootstrap.Run().Env.PrimaryRedis)
	cr := redis.NewUserCacheRepository(rdb)
	m := mail.NewEmailService()
	s := service.NewGeneralService(uw, cr, m)
	uc := general.NewGeneralUserController(s)

	auth := v1.Group("/auth")
	{
		auth.POST("/login", uc.Login)
		auth.POST("/register", uc.Register)
		// auth.POST("/verify", uc.VerifyEmail)
		// auth.POST("/forgotpassword", uc.ForgotPassword)
	}
}
