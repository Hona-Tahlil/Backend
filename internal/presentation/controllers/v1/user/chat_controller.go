package user

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
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
	client := websocket.NewClient(c.hub, conn, params.RoomID, userID, &bootstrap.Run().Env.WebsocketSetting, c.chatService)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func (c *UserChatController) GetAllRooms(ctx *gin.Context) {
	type Sort struct {
		Field string `form:"field" validate:"required"`
		Dir   string `form:"dir"   validate:"required,oneof=ASC DESC"`
	}
	type Params struct {
		Page  int    `form:"page"`
		Count int    `form:"count"`
		Sort  []Sort `form:"sort"`
	}
	senderID := controllers.GetID(ctx)
	p := controllers.Receive[Params](ctx)
	offset, limit := controllers.GetOffsetLimit(p.Page, p.Count, 1, 10)

	sorts := make([]general.Sort, len(p.Sort))
	for i, s := range p.Sort {
		sorts[i] = general.Sort{
			Field: s.Field,
			Dir:   s.Dir,
		}
	}

	request := chat.GetAllRoomsRequest{
		SenderID: senderID,
		Offset:   offset,
		Limit:    limit,
		Sort:     sorts,
	}
	rooms, totalCount, err := c.chatService.GetAllRooms(request)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(rooms, totalCount, offset, limit)
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}

func (c *UserChatController) GetRoomMessages(ctx *gin.Context) {
	type Sort struct {
		Field string `form:"field" validate:"required"`
		Dir   string `form:"dir"   validate:"required,oneof=ASC DESC"`
	}
	type params struct {
		RoomID uint   `uri:"roomID" validate:"required"`
		page   int    `form:"page" validate:"required"`
		count  int    `form:"count" validate:"required"`
		Sort   []Sort `form:"sort" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	offset, limit := controllers.GetOffsetLimit(p.page, p.count, 1, 10)

	sorts := make([]general.Sort, len(p.Sort))
	for i, s := range p.Sort {
		sorts[i] = general.Sort{
			Field: s.Field,
			Dir:   s.Dir,
		}
	}
	request := chat.GetRoomMessagesRequest{
		RoomID:   p.RoomID,
		SenderID: senderID,
		Offset:   offset,
		Limit:    limit,
		Sort:     sorts,
	}
	messages, totalCount, err := c.chatService.GetRoomMessages(request)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(messages, totalCount, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}


func (c *UserChatController) BlockRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	request := chat.BlockRoomRequest{
		RoomID:    p.RoomID,
		SenderID:  senderID,
		BlockedBy: enums.C_BlockedByUser,
	}
	if err := c.chatService.BlockRoom(request); err != nil {
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
	senderID := controllers.GetID(ctx)
	request := chat.UnblockRoomRequest{
		RoomID:   p.RoomID,
		SenderID: senderID,
	}
	if err := c.chatService.UnblockRoom(request); err != nil {
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

