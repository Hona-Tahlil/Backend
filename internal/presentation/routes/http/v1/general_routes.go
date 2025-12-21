package httpv1

import (
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
		auth.PUT("/reset-password", app.Controllers.GeneralControllers.GeneralUserController.ResetPassword)
	}

	provinces := v1.Group("/provinces")
	{
		provinces.GET("/", app.Controllers.GeneralControllers.GeneralProvinceController.GetAllProvinces)
		provinces.GET("/:province_num/cities", app.Controllers.GeneralControllers.GeneralProvinceController.GetCitiesByProvinceName)
	}
	pets := v1.Group("/pets")
	{
		pets.GET("/kinds", app.Controllers.GeneralControllers.GeneralPetController.GetAllPetKinds)
		pets.GET("/kinds/:petKind/species", app.Controllers.GeneralControllers.GeneralPetController.GetPetKindSpecies)
	}
	search := v1.Group("/search")
	{
		search.GET("/petsitter", app.Controllers.GeneralControllers.GeneralSearchController.SearchPetSitters)
	}
}
