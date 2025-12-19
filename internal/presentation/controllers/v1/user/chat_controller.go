package user

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type UserChatController struct {
	chatService usecase.ChatService
	hub         *websocket.Hub
}

func NewUserChatController(chatService usecase.ChatService, hub *websocket.Hub) *UserChatController {
	return &UserChatController{
		chatService: chatService,
		hub:         hub,
	}
}

func (c *UserChatController) CreateOrGetRoom(ctx *gin.Context) {
	// Implementation for creating or getting a chat room
	type roomParams struct {
		PetSitterID uint `uri:"petSitterID" validate:"required"`
	}
	params := controllers.Receive[roomParams](ctx)
	userID := controllers.GetID(ctx)
	roomInfo := chat.CreateOrGetUserRoomRequest{
		PetSitterID: params.PetSitterID,
		UserID:      userID,
	}
	roomsDetails, err := c.chatService.CreateOrGetRoom(roomInfo)
	if err != nil {
		panic(err)
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, roomsDetails)
}

func (c *UserChatController) HandleWebsocket(ctx *gin.Context) {
	type wsParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	params := controllers.Receive[wsParams](ctx)
	userID := controllers.GetID(ctx)
	conn, _ := ctx.Get(bootstrap.Run().Constants.Context.WebsocketConnection)
	// Implementation for handling websocket connection
	// client := websocket.NewClient(chatController.hub, conn, params.RoomID, userID, chatController.websocketSetting, chatController.chatService, nil)
	//new client
	client := websocket.NewClient(c.hub, conn, params.RoomID, userID, &bootstrap.Run().Env.WebsocketSetting, c.chatService)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func (c *UserChatController) GetAllRooms(ctx *gin.Context) {
	// For now return empty list or delegate to service when available
	// Query params can be added later (e.g., status)
	// userID := controllers.GetID(ctx)
	// res, err := c.chatService.GetAllRooms(userID)
	// if err != nil { panic(err) }
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, []interface{}{})
}

func (c *UserChatController) BlockRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	userID := controllers.GetID(ctx)

	if err := c.chatService.BlockRoom(userID, p.RoomID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *UserChatController) UnblockRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	userID := controllers.GetID(ctx)

	if err := c.chatService.UnblockRoom(userID, p.RoomID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *UserChatController) GetRoomRequestInfo(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)

	res, err := c.chatService.GetRoomRequestInfo(p.RoomID)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
