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
		requests.GET("/", app.Controllers.UserControllers.UserRequestController.GetCreateRequestInfo)
		requests.POST("/", app.Controllers.UserControllers.UserRequestController.CreateRequest)
		requests.POST("/search", app.Controllers.UserControllers.UserRequestController.SearchRequests)
		requests.PUT("/", app.Controllers.UserControllers.UserRequestController.EditRequest)
		requests.PUT("/cancel", app.Controllers.UserControllers.UserRequestController.CancelRequest)
		requests.PUT("/pay", app.Controllers.UserControllers.UserRequestController.PayRequest)
		requests.GET("/:requestID", app.Controllers.UserControllers.UserRequestController.GetRequestFullData)
	}

	comments := v1.Group("/comments")
	comments.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		comments.POST("/", app.Controllers.UserControllers.UserCommentController.CreateComment)
		comments.PUT("/", app.Controllers.UserControllers.UserCommentController.EditComment)
		comments.DELETE("/:id", app.Controllers.UserControllers.UserCommentController.DeleteComment)
		comments.GET("/petsitters/:petSitterID", app.Controllers.UserControllers.UserCommentController.GetAllPetSitterComments)
	}
}
