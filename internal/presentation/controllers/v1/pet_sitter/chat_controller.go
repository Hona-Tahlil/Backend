package petsitter

import (
	"encoding/json"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/infrastructure/jwt"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/controllers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

type PetSitterChatController struct {
	chatService usecase.ChatService
	hub         *websocket.Hub
	jwtService  *jwt.JWTService
}

func NewPetSitterChatController(chatService usecase.ChatService, hub *websocket.Hub, jwtService *jwt.JWTService) *PetSitterChatController {
	return &PetSitterChatController{
		chatService: chatService,
		hub:         hub,
		jwtService:  jwtService,
	}
}

// func (c *PetSitterChatController) HandleWebsocket(ctx *gin.Context) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("WebSocket handler panic recovered:", r)
// 		}
// 	}()

// 	type wsParams struct {
// 		RoomID uint `uri:"roomID" validate:"required"`
// 	}

// 	params := controllers.Receive[wsParams](ctx)
// 	if params.RoomID == 0 {
// 		log.Println("Invalid RoomID in WebSocket request")
// 		return
// 	}

// 	userID := controllers.GetID(ctx)
// 	if userID == 0 {
// 		log.Println("Invalid UserID in WebSocket request")
// 		return
// 	}

// 	conn, exists := ctx.Get(bootstrap.Run().Constants.Context.WebsocketConnection)
// 	if !exists || conn == nil {
// 		log.Println("WebSocket connection not found in context")
// 		return
// 	}

// 	client := websocket.NewClient(c.hub, conn, params.RoomID, userID, &bootstrap.Run().Env.WebsocketSetting, c.chatService)
// 	if client == nil {
// 		log.Println("Failed to create WebSocket client")
// 		return
// 	}

// 	log.Println("WebSocket client created successfully for room:", params.RoomID)
// 	client.Hub.Register <- client

//		go client.ReadPump()
//		go client.WritePump()
//	}
func (c *PetSitterChatController) HandleWebsocket(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("WebSocket handler panic recovered:", r)
		}
	}()

	type wsParams struct {
		RoomID uint   `uri:"roomID" validate:"required"`
		Token  string `uri:"token" validate:"required"`
	}

	params := controllers.Receive[wsParams](ctx)
	if params.RoomID == 0 {
		log.Println("Invalid RoomID")
		return
	}

	// userID := controllers.GetID(ctx)
	// if userID == 0 {
	// 	log.Println("Invalid UserID")
	// 	return
	// }
	userID := c.jwtService.ValidateToken(params.Token, bootstrap.Run().Constants.JWTConstants.AccessTokenType)

	// claims, err := c.jwtService.ValidateToken(params.Token)
	// if err != nil {
	// 	panic(err)
	// }
	// userID := uint(claims["sub"].(float64))
	upgrader := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := websocket.NewClient(c.hub, wsConn, params.RoomID, userID, &bootstrap.Run().Env.WebsocketSetting, c.chatService)
	if client == nil {
		wsConn.Close()
		log.Println("Failed to create WebSocket client")
		return
	}

	client.Hub.Register <- client
	log.Println("WebSocket client created successfully for room:", params.RoomID)
	go client.ReadPump()
	go client.WritePump()
}

func (c *PetSitterChatController) GetAllRooms(ctx *gin.Context) {
	type Params struct {
		Page  int `form:"page"`
		Count int `form:"count"`
	}
	senderID := controllers.GetID(ctx)
	p := controllers.Receive[Params](ctx)
	offset, limit := controllers.GetOffsetLimit(p.Page, p.Count)

	request := chat.GetAllRoomsRequest{
		SenderID: senderID,
		Offset:   offset,
		Limit:    limit,
	}
	rooms, totalCount, err := c.chatService.GetAllRooms(request)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(rooms, totalCount, offset, limit)
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}

func (c *PetSitterChatController) GetRoomMessages(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
		page   int  `form:"page" validate:"required"`
		count  int  `form:"count" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	offset, limit := controllers.GetOffsetLimit(p.page, p.count)

	request := &chat.GetRoomMessagesRequest{
		RoomID:   p.RoomID,
		SenderID: senderID,
		Offset:   offset,
		Limit:    limit,
	}
	messages, totalCount, err := c.chatService.GetRoomMessages(request)
	if err != nil {
		panic(err)
	}
	data := controllers.NewPaginatedResponse(messages, totalCount, offset, limit)

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, data)
}

func (c *PetSitterChatController) EditMessage(ctx *gin.Context) {
	type params struct {
		RoomID      uint                 `uri:"roomID" validate:"required"`
		MessageID   uint                 `uri:"messageID" validate:"required"`
		Content     string               `json:"content" form:"content"`
		MessageType enums.ChatMessageType `json:"message_type" form:"message_type"`
		MediaBase64 string               `json:"media_base64" form:"media_base64"`
		MediaMime   string               `json:"media_mime" form:"media_mime"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	request := chat.EditMessageRequest{
		RoomID:      p.RoomID,
		MessageID:   p.MessageID,
		SenderID:    senderID,
		Content:     p.Content,
		MessageType: p.MessageType,
		MediaBase64: p.MediaBase64,
		MediaMime:   p.MediaMime,
	}
	res, err := c.chatService.EditMessage(request)
	if err != nil {
		panic(err)
	}
	payload := websocket.EditPayload{
		MessageID:   res.ID,
		MessageType: string(res.MessageType),
		Content:     res.Content,
		MediaURL:    res.MediaURL,
		IsEdited:    res.IsEdited,
		EditedAt:    time.Now().UTC(),
	}
	if b, err := json.Marshal(payload); err == nil {
		c.hub.Broadcast <- &websocket.Message{
			Type:      websocket.MessageTypeEdit,
			RoomID:    p.RoomID,
			SenderID:  senderID,
			MessageID: res.ID,
			Timestamp: time.Now().UTC(),
			Content:   b,
		}
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (c *PetSitterChatController) DeleteMessage(ctx *gin.Context) {
	type params struct {
		RoomID    uint `uri:"roomID" validate:"required"`
		MessageID uint `uri:"messageID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	request := chat.DeleteMessageRequest{
		RoomID:    p.RoomID,
		MessageID: p.MessageID,
		SenderID:  senderID,
	}
	res, err := c.chatService.DeleteMessage(request)
	if err != nil {
		panic(err)
	}
	payload := websocket.DeletePayload{
		MessageID: res.ID,
		IsDeleted: res.IsDeleted,
		DeletedAt: res.DeletedAt,
	}
	if b, err := json.Marshal(payload); err == nil {
		c.hub.Broadcast <- &websocket.Message{
			Type:      websocket.MessageTypeDelete,
			RoomID:    p.RoomID,
			SenderID:  senderID,
			MessageID: res.ID,
			Timestamp: time.Now().UTC(),
			Content:   b,
		}
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}

func (c *PetSitterChatController) AddReaction(ctx *gin.Context) {
	type params struct {
		RoomID    uint   `uri:"roomID" validate:"required"`
		MessageID uint   `uri:"messageID" validate:"required"`
		Emoji     string `json:"emoji" form:"emoji" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	request := chat.ReactionRequest{
		RoomID:    p.RoomID,
		MessageID: p.MessageID,
		SenderID:  senderID,
		Emoji:     p.Emoji,
	}
	reactions, err := c.chatService.AddReaction(request)
	if err != nil {
		panic(err)
	}
	payload := websocket.ReactionPayload{
		MessageID: p.MessageID,
		Emoji:     p.Emoji,
		Action:    "ADD",
		UserID:    senderID,
		Timestamp: time.Now().UTC(),
	}
	if b, err := json.Marshal(payload); err == nil {
		c.hub.Broadcast <- &websocket.Message{
			Type:      websocket.MessageTypeReaction,
			RoomID:    p.RoomID,
			SenderID:  senderID,
			MessageID: p.MessageID,
			Timestamp: payload.Timestamp,
			Content:   b,
		}
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, reactions)
}

func (c *PetSitterChatController) RemoveReaction(ctx *gin.Context) {
	type params struct {
		RoomID    uint   `uri:"roomID" validate:"required"`
		MessageID uint   `uri:"messageID" validate:"required"`
		Emoji     string `json:"emoji" form:"emoji" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)
	request := chat.ReactionRequest{
		RoomID:    p.RoomID,
		MessageID: p.MessageID,
		SenderID:  senderID,
		Emoji:     p.Emoji,
	}
	reactions, err := c.chatService.RemoveReaction(request)
	if err != nil {
		panic(err)
	}
	payload := websocket.ReactionPayload{
		MessageID: p.MessageID,
		Emoji:     p.Emoji,
		Action:    "REMOVE",
		UserID:    senderID,
		Timestamp: time.Now().UTC(),
	}
	if b, err := json.Marshal(payload); err == nil {
		c.hub.Broadcast <- &websocket.Message{
			Type:      websocket.MessageTypeReaction,
			RoomID:    p.RoomID,
			SenderID:  senderID,
			MessageID: p.MessageID,
			Timestamp: payload.Timestamp,
			Content:   b,
		}
	}
	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, reactions)
}

func (c *PetSitterChatController) AcceptRoom(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)

	if err := c.chatService.AcceptRoom(senderID, p.RoomID); err != nil {
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
	senderID := controllers.GetID(ctx)

	if err := c.chatService.RejectRoom(senderID, p.RoomID); err != nil {
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
	senderID := controllers.GetID(ctx)
	request := chat.BlockRoomRequest{
		RoomID:    p.RoomID,
		SenderID:  senderID,
		BlockedBy: enums.C_BlockedByPetSitter,
	}
	if err := c.chatService.BlockRoom(request); err != nil {
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

func (c *PetSitterChatController) GetRoomRequestInfo(ctx *gin.Context) {
	type params struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	p := controllers.Receive[params](ctx)
	senderID := controllers.GetID(ctx)

	res, err := c.chatService.GetRoomRequestInfo(p.RoomID, senderID)
	if err != nil {
		panic(err)
	}

	msg := controllers.Message{}
	controllers.Respond(ctx, 200, msg, res)
}
