package httpv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	requests := v1.Group("/requests")
	requests.Use(app.Middlewares.AuthMiddleware.AuthRequired)
	{
		requests.GET("/:requestID", app.Controllers.PetSitterControllers.PetSitterRequestController.GetRequestFullData)
		requests.PUT("/cancel", app.Controllers.PetSitterControllers.PetSitterRequestController.CancelRequest)
		requests.PUT("/respond", app.Controllers.PetSitterControllers.PetSitterRequestController.RespondToRequest)
	}
}
