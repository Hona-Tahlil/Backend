package httpv1

import (
	"hona/backend/internal/presentation/controllers/v1/user"
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login", app.Controllers.GeneralControllers.GeneralUserController.Login)
		auth.POST("/register", app.Controllers.GeneralControllers.GeneralUserController.Register)
		auth.POST("/verify", app.Controllers.GeneralControllers.GeneralUserController.VerifyEmail)
		auth.POST("/forgot-password", app.Controllers.GeneralControllers.GeneralUserController.ForgotPassword)
	}
	v1.POST("/file", user.NewUserPetController().AddPet)
}
