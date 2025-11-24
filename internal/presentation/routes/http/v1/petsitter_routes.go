package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	petSitterRequests := v1.Group("/petsitter-requests")
	petSitterRequests.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		// TODO: test
		petSitterRequests.GET("/:requestID", app.Controllers.PetSitterControllers.PetSitterRequestController.GetRequestFullData)
		// TODO: test
		petSitterRequests.PUT("/cancel", app.Controllers.PetSitterControllers.PetSitterRequestController.CancelRequest)
		// TODO: test
		petSitterRequests.PUT("/respond", app.Controllers.PetSitterControllers.PetSitterRequestController.RespondToRequest)
		// TODO: add View Requests Routes
	}
}
