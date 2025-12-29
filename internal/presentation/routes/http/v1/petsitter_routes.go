package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	petsitter := v1.Group("/petsitter")
	petsitter.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		register := petsitter.Group("/register")
		{
			// STEP 0: Create signup session
			register.POST("/", app.Controllers.PetSitterControllers.PetSitterRegisterController.CreateSignupSession)

			// STEP 1: Personal info
			register.PUT("/personal", app.Controllers.PetSitterControllers.PetSitterRegisterController.SubmitPersonalInfo)
			register.GET("/personal", app.Controllers.PetSitterControllers.PetSitterRegisterController.GetPersonalInfo)

			// STEP 2: Upload documents
			register.PUT("/documents", app.Controllers.PetSitterControllers.PetSitterRegisterController.UploadDocuments)
			register.GET("/documents", app.Controllers.PetSitterControllers.PetSitterRegisterController.GetDocuments)

			// STEP 3: Skills + Bio
			register.PUT("/skills", app.Controllers.PetSitterControllers.PetSitterRegisterController.SubmitSkills)
			// register.GET("/skills", app.Controllers.PetSitterController.PetSitterController.GetSkills)

			// Optional: Get current status
			register.GET("/status", app.Controllers.PetSitterControllers.PetSitterRegisterController.GetPetsitterStatus)
		}
		requests := petsitter.Group("/requests")
		{
			requests.POST("/search", app.Controllers.PetSitterControllers.PetSitterRequestController.SearchRequests)
			requests.GET("/:requestID", app.Controllers.PetSitterControllers.PetSitterRequestController.GetRequestFullData)
			requests.PUT("/cancel", app.Controllers.PetSitterControllers.PetSitterRequestController.CancelRequest)
			requests.PUT("/respond", app.Controllers.PetSitterControllers.PetSitterRequestController.RespondToRequest)
		}
		chat := petsitter.Group("/chat")
		{
			//get by status
			chat.POST("/room/:userID", app.Controllers.UserControllers.UserChatController.CreateOrGetRoom)
			chat.GET("/rooms", app.Controllers.PetSitterControllers.PetSitterChatController.GetAllRooms)
			chat.PUT("/room/:roomID/accept", app.Controllers.PetSitterControllers.PetSitterChatController.AcceptRoom)
			chat.PUT("/room/:roomID/reject", app.Controllers.PetSitterControllers.PetSitterChatController.RejectRoom)
			chat.PUT("/room/:roomID/block", app.Controllers.PetSitterControllers.PetSitterChatController.BlockRoom)
			chat.PUT("/room/:roomID/unblock", app.Controllers.PetSitterControllers.PetSitterChatController.UnblockRoom)
			chat.GET("/room/:roomID/request-info", app.Controllers.PetSitterControllers.PetSitterChatController.GetRoomRequestInfo)
			chat.GET("/room/:roomID/messages", app.Controllers.PetSitterControllers.PetSitterChatController.GetRoomMessages)

		}
	}

	petKinds := petsitter.Group("/pet-kinds")
	{
		petKinds.GET("/", app.Controllers.PetSitterControllers.PetSitterSkillsController.GetPetKinds)
		petKinds.PUT("/", app.Controllers.PetSitterControllers.PetSitterSkillsController.UpdatePetKinds)
	}

	services := petsitter.Group("/services")
	{
		services.GET("/", app.Controllers.PetSitterControllers.PetSitterSkillsController.GetServices)
		services.POST("/", app.Controllers.PetSitterControllers.PetSitterSkillsController.CreateService)
		services.PUT("/", app.Controllers.PetSitterControllers.PetSitterSkillsController.UpdateService)
		services.DELETE("/:serviceID", app.Controllers.PetSitterControllers.PetSitterSkillsController.DeleteService)
	}

	calendar := petsitter.Group("/calendar")
	{
		calendar.GET("/", app.Controllers.PetSitterControllers.PetSitterCalendarController.GetCalendarSlots)
		calendar.PATCH("/", app.Controllers.PetSitterControllers.PetSitterCalendarController.UpdateFreeCalendarSlots)
	}

	wallet := petsitter.Group("/wallet")
	{
		wallet.PUT("/withdraw", app.Controllers.PetSitterControllers.PetSitterWalletController.Withdraw)
	}
}
