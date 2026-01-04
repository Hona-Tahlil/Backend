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
	chat := v1.Group("/chat")
	chat.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		chat.POST("/rooms/:petSitterID", app.Controllers.UserControllers.UserChatController.CreateOrGetRoom)
		chat.GET("/rooms", app.Controllers.UserControllers.UserChatController.GetAllRooms)
		chat.PUT("/room/:roomID/block", app.Controllers.UserControllers.UserChatController.BlockRoom)
		chat.PUT("/room/:roomID/unblock", app.Controllers.UserControllers.UserChatController.UnblockRoom)
		chat.GET("/room/:roomID/request-info", app.Controllers.UserControllers.UserChatController.GetRoomRequestInfo)
		chat.GET("/room/:roomID/messages", app.Controllers.UserControllers.UserChatController.GetRoomMessages)
		chat.PUT("/room/:roomID/message/:messageID", app.Controllers.UserControllers.UserChatController.EditMessage)
		chat.DELETE("/room/:roomID/message/:messageID", app.Controllers.UserControllers.UserChatController.DeleteMessage)
		chat.POST("/room/:roomID/message/:messageID/reaction", app.Controllers.UserControllers.UserChatController.AddReaction)
		chat.DELETE("/room/:roomID/message/:messageID/reaction", app.Controllers.UserControllers.UserChatController.RemoveReaction)
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
