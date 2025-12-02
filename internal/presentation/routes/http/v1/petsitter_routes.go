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
			register.POST("/", app.Controllers.PetSitterControllers.PetSitterController.CreateSignupSession)


			// STEP 1: Personal info
			register.PUT("/personal", app.Controllers.PetSitterControllers.PetSitterController.SubmitPersonalInfo)
			register.GET("/personal", app.Controllers.PetSitterControllers.PetSitterController.GetPersonalInfo)


			// STEP 2: Upload documents
			register.PUT("/documents", app.Controllers.PetSitterControllers.PetSitterController.UploadDocuments)
			register.GET("/documents", app.Controllers.PetSitterControllers.PetSitterController.GetDocuments)


			// STEP 3: Skills + Bio
			register.PUT("/skills", app.Controllers.PetSitterControllers.PetSitterController.SubmitSkills)
			// register.GET("/skills", app.Controllers.PetSitterController.PetSitterController.GetSkills)


			// Optional: Get current status
			register.GET("/status", app.Controllers.PetSitterControllers.PetSitterController.GetPetsitterStatus)
		}
		requests := petsitter.Group("/requests")
		{
			// TODO: test
			requests.GET("/:requestID", app.Controllers.PetSitterControllers.PetSitterRequestController.GetRequestFullData)
			// TODO: test
			requests.PUT("/cancel", app.Controllers.PetSitterControllers.PetSitterRequestController.CancelRequest)
			// TODO: test
			requests.PUT("/respond", app.Controllers.PetSitterControllers.PetSitterRequestController.RespondToRequest)
			// TODO: add View Requests Routes
		}
	}
}
