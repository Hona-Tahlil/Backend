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
	RoomID  uint
	SenderID  uint
	Content string
	ReplyToMessageID *uint
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
	Sort     []general.Sort
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