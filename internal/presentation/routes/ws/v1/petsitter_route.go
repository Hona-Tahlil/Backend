package wsv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	chat := v1.Group("/petsitter/chat")
	// chat.Use(app.Middlewares.WebsocketMiddleware.UpgradeToWebSocket)
	{
		chat.GET("/room/:roomID/:token", app.Controllers.PetSitterControllers.PetSitterChatController.HandleWebsocket)
	}
}
