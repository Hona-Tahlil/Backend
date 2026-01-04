package chat

import (
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/domain/enums"
)

type CreateOrGetUserRoomRequest struct {
	PetSitterID uint
	UserID      uint
}

type SaveMessageRequest struct {
	RoomID           uint
	SenderID         uint
	Content          string
	MessageType      enums.ChatMessageType
	MediaBase64      string
	MediaMime        string
	ReplyToMessageID *uint
}

type EditMessageRequest struct {
	RoomID      uint
	MessageID   uint
	SenderID    uint
	Content     string
	MessageType enums.ChatMessageType
	MediaBase64 string
	MediaMime   string
}

type DeleteMessageRequest struct {
	RoomID    uint
	MessageID uint
	SenderID  uint
}

type ReactionRequest struct {
	RoomID    uint
	MessageID uint
	SenderID  uint
	Emoji     string
}

type BlockRoomRequest struct {
	RoomID    uint
	SenderID  uint
	BlockedBy enums.ChatBlockedBy
}

type UnblockRoomRequest struct {
	RoomID   uint
	SenderID uint
}

type GetAllRoomsRequest struct {
	SenderID uint
	Offset   int
	Limit    int
}

type GetRoomMessagesRequest struct {
	RoomID uint
	SenderID uint
	Offset int
	Limit  int
	Sort   []general.Sort
}

type MarkRoomReadRequest struct {
	RoomID uint
	SenderID uint
	LastReadMessageID uint
}
