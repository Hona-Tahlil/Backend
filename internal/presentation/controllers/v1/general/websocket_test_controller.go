package general

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type GeneralWebsocketTestController struct {
	hub *websocket.Hub
}

func NewGeneralWebsocketTestController(hub *websocket.Hub) *GeneralWebsocketTestController {
	return &GeneralWebsocketTestController{
		hub: hub,
	}
}

// HandleWebsocketTest handles WebSocket connection for testing purposes
// This endpoint allows testing WebSocket functionality without authentication
func (c *GeneralWebsocketTestController) HandleWebsocketTest(ctx *gin.Context) {
	type wsParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	params := controllers.Receive[wsParams](ctx)

	// For testing, use a fixed test user ID
	testUserID := uint(999)

	conn, _ := ctx.Get(bootstrap.Run().Constants.Context.WebsocketConnection)

	// Create a new client for testing
	client := websocket.NewClient(c.hub, conn, params.RoomID, testUserID, &bootstrap.Run().Env.WebsocketSetting, nil)
	client.Hub.Register <- client
	go client.ReadPump()
	go client.WritePump()
}

// GetWebsocketStatus returns the current WebSocket hub status (for testing/debugging)
func (c *GeneralWebsocketTestController) GetWebsocketStatus(ctx *gin.Context) {
	type StatusResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	msg := controllers.Message{
		Text: "WebSocket hub is running",
	}

	controllers.Respond(ctx, 200, msg, StatusResponse{
		Message: "WebSocket testing endpoint is active",
		Status:  "ok",
	})
}
