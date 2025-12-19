package petsitter

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/controllers"

	"github.com/gin-gonic/gin"
)

type PetSitterChatController struct {
	chatService usecase.ChatService
	hub         *websocket.Hub
}

func NewPetSitterChatController(chatService usecase.ChatService, hub *websocket.Hub) *PetSitterChatController {
	return &PetSitterChatController{
		chatService: chatService,
		hub:         hub,
	}
}

func (c *PetSitterChatController) HandleWebsocket(ctx *gin.Context) {
	// Implementation for handling websocket connection
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

func (c *PetSitterChatController) GetAllRooms(ctx *gin.Context) {
	// For now return empty list or delegate to service when available
	// Query params can be added later (e.g., status)
	// petSitterID := controllers.GetID(ctx)
	// res, err := c.chatService.GetAllRooms(petSitterID)
	// if err != nil { panic(err) }
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, []interface{}{})
}

func (c *PetSitterChatController) AcceptRoom(ctx *gin.Context) {
	type params struct {
		UserID uint `uri:"userID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	petSitterID := controllers.GetID(ctx)

	if err := c.chatService.AcceptRoom(petSitterID, p.UserID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *PetSitterChatController) RejectRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	petSitterID := controllers.GetID(ctx)
	_ = petSitterID

	if err := c.chatService.RejectRoom(petSitterID, p.RoomID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *PetSitterChatController) BlockRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	petSitterID := controllers.GetID(ctx)

	if err := c.chatService.BlockRoom(petSitterID, p.RoomID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *PetSitterChatController) UnblockRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	petSitterID := controllers.GetID(ctx)

	if err := c.chatService.UnblockRoom(petSitterID, p.RoomID); err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, nil)
}

func (c *PetSitterChatController) GetRoomRequestInfo(ctx *gin.Context) {
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
