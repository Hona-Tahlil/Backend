package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpUserRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	pets := v1.Group("/pets")
	pets.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		pets.POST("/", app.Controllers.UserControllers.UserPetController.AddPet)
		pets.PUT("/", app.Controllers.UserControllers.UserPetController.UpdatePet)
		pets.DELETE("/:id", app.Controllers.UserControllers.UserPetController.RemovePet)
		pets.GET("/", app.Controllers.UserControllers.UserPetController.GetPetsBasicData)
		pets.GET("/:id", app.Controllers.UserControllers.UserPetController.GetPetFullData)
	}

	requests := v1.Group("/requests")
	requests.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		// TODO: test
		requests.GET("/create-info", app.Controllers.UserControllers.UserRequestController.GetCreateRequestInfo)
		// TODO: test
		requests.POST("/", app.Controllers.UserControllers.UserRequestController.CreateRequest)
		// TODO: test
		requests.PUT("/", app.Controllers.UserControllers.UserRequestController.EditRequest)
		// TODO: test
		requests.PUT("/cancel", app.Controllers.UserControllers.UserRequestController.CancelRequest)
		// TODO: test
		requests.GET("/:requestID", app.Controllers.UserControllers.UserRequestController.GetRequestFullData)
		// TODO: add View Requests Routes
	}
}
