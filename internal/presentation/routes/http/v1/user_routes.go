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
		comments.DELETE("/:commentID", app.Controllers.UserControllers.UserCommentController.DeleteComment)
		comments.GET("/petsitters/:petSitterID", app.Controllers.UserControllers.UserCommentController.GetAllPetSitterComments)
	}

	profile := v1.Group("/profile")
	profile.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		profile.GET("/", app.Controllers.UserControllers.UserProfileController.GetProfile)
		profile.GET("/identity", app.Controllers.UserControllers.UserProfileController.GetIdentity)
		profile.PUT("/", app.Controllers.UserControllers.UserProfileController.UpdateProfile)
	}

	wallet := v1.Group("/wallet")
	wallet.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		wallet.GET("/", app.Controllers.UserControllers.UserWalletController.GetWallet)
		wallet.GET("/transfers", app.Controllers.UserControllers.UserWalletController.ListTransfers)
		wallet.GET("/transactions", app.Controllers.UserControllers.UserWalletController.ListTransactions)
		wallet.PUT("/top-up", app.Controllers.UserControllers.UserWalletController.TopUp)
	}
}
