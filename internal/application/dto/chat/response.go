package chat

import (
	"hona/backend/internal/domain/enums"
	"time"
)

type ChatRoomDetailsResponse struct {
	RoomID    uint                    `json:"room_id"`
	User      ChatParticipantResponse `json:"user"`
	PetSitter ChatParticipantResponse `json:"pet_sitter"`
	Status    string                  `json:"status"` // PENDING | ACCEPTED | REJECTED | BLOCKED
	// IsAccepted bool                    `json:"is_accepted"`
	BlockedBy *string `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
}

type ChatParticipantResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
	IsOnline   bool   `json:"is_online"`
}

type RoomMessagesResponse struct {
	ID               uint                    `json:"id"`
	RoomID           uint                    `json:"room_id"`
	Sender           ChatParticipantResponse `json:"sender"`
	Content          string                  `json:"content"`
	CreatedAt        time.Time               `json:"created_at"`
	IsUnread         bool                    `json:"is_unread"`
	IsMine           bool                    `json:"is_mine"`
	ReplyToMessageID *uint                   `json:"reply_to_message_id,omitempty"`
	// IsEdited  bool                    `json:"is_edited"`
}

type UserRoomsResponse struct {
	RoomID             uint                    `json:"room_id"`
	PetSitter          ChatParticipantResponse `json:"pet_sitter"`
	Status             enums.ChatRoomStatus    `json:"status"`               // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy          *enums.ChatBlockedBy    `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
	LastMessage        *string                 `json:"last_message,omitempty"`
	LastMessageTime    *time.Time              `json:"last_message_time,omitempty"`
	UnreadMessageCount uint                    `json:"unread_message_count"`
}

type PetSitterRoomsResponse struct {
	RoomID             uint                    `json:"room_id"`
	User               ChatParticipantResponse `json:"user"`
	Status             enums.ChatRoomStatus    `json:"status"`               // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy          *enums.ChatBlockedBy    `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
	LastMessage        *string                 `json:"last_message,omitempty"`
	LastMessageTime    *time.Time              `json:"last_message_time,omitempty"`
	UnreadMessageCount int64                   `json:"unread_message_count"`
}

type SaveMessageResponse struct {
	ID        uint      `json:"id"`
	RoomID    uint      `json:"room_id"`	
	SenderID  uint      `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
