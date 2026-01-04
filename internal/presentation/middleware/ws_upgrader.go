package middleware

import (
	"hona/backend/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMiddleware struct {
	constants *bootstrap.Constants
}

func NewWebsocketMiddleware(constants *bootstrap.Constants) *WebsocketMiddleware {
	return &WebsocketMiddleware{
		constants: constants,
	}
}

func (wsMiddleware *WebsocketMiddleware) UpgradeToWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Set(wsMiddleware.constants.Context.WebsocketConnection, conn)
	c.Abort()
}

// func (wsMiddleware *WebsocketMiddleware) UpgradeToWebSocket(c *gin.Context) {
// 	upgrader := websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool { return true },
// 	}
// 	log.Println("Upgrading connection to WebSocket")
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	c.Set(wsMiddleware.constants.Context.WebsocketConnection, conn)

// 	c.Abort()
// }
