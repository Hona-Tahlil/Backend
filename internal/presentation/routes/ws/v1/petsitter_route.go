package wsv1

import (
	"hona/backend/wire"

	"github.com/gin-gonic/gin"
)

func SetUpPetSitterRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	chat := v1.Group("/chat")
	{
		chat.GET("/room/:roomID", app.Controllers.PetSitterControllers.PetSitterChatController.HandleWebsocket)
	}

}
