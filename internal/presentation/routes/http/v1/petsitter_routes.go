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
	}
}
