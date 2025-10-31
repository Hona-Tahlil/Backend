package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpUserRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	pets := v1.Group("/pets")
	{
		pets.POST("/", app.Controllers.UserControllers.UserPetController.AddPet)
		pets.PUT("/", app.Controllers.UserControllers.UserPetController.UpdatePet)
		pets.DELETE("/:id", app.Controllers.UserControllers.UserPetController.RemovePet)
		pets.GET("/", app.Controllers.UserControllers.UserPetController.GetPetsBasicData)
		pets.GET("/:id", app.Controllers.UserControllers.UserPetController.GetPetFullData)
	}
}
