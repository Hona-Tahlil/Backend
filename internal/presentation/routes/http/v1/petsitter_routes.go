package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	petsitter := v1.Group("/petsitter")
	{
		register := petsitter.Group("/register")
		{
			// register.POST("/status", app.Controllers.PetSitterController.PetSitterController.GetStatus)
			// register.GET("/step1", app.Controllers.PetSitterController.PetSitterController.Register)
			// register.PUT("/step1", app.Controllers.PetSitterController.PetSitterController.Register)
			// register.GET("/step2", app.Controllers.PetSitterController.PetSitterController.Register)
			// register.PUT("/step2", app.Controllers.PetSitterController.PetSitterController.Register)
			// register.GET("/step3", app.Controllers.PetSitterController.PetSitterController.Register)
			// register.PUT("/step3", app.Controllers.PetSitterController.PetSitterController.Register)
			// STEP 0: Create signup session
			register.POST("/", app.Controllers.PetSitterController.PetSitterController.CreateSignupSession)


			// STEP 1: Personal info
			register.PUT("/personal", app.Controllers.PetSitterController.PetSitterController.SubmitPersonalInfo)
			register.GET("/personal", app.Controllers.PetSitterController.PetSitterController.GetPersonalInfo)


			// STEP 2: Upload documents
			register.PUT("/documents", app.Controllers.PetSitterController.PetSitterController.UploadDocuments)
			register.GET("/documents", app.Controllers.PetSitterController.PetSitterController.GetDocuments)


			// STEP 3: Skills + Bio
			register.PUT("/skills", app.Controllers.PetSitterController.PetSitterController.SubmitSkills)
			register.GET("/skills", app.Controllers.PetSitterController.PetSitterController.GetSkills)


			// Optional: Get current status
			register.GET("/status", app.Controllers.PetSitterController.PetSitterController.GetSignupStatus)
		}
	}
}
