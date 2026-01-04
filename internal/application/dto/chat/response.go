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
	LastSeen   *time.Time `json:"last_seen,omitempty"`
}

type RoomMessagesResponse struct {
	ID               uint                    `json:"id"`
	RoomID           uint                    `json:"room_id"`
	Sender           ChatParticipantResponse `json:"sender"`
	MessageType      enums.ChatMessageType   `json:"message_type"`
	Content          string                  `json:"content"`
	MediaURL         *string                 `json:"media_url,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	IsEdited         bool                    `json:"is_edited"`
	IsDeleted        bool                    `json:"is_deleted"`
	DeletedAt        *time.Time              `json:"deleted_at,omitempty"`
	Reactions        []MessageReactionSummary `json:"reactions,omitempty"`
	IsUnread         bool                    `json:"is_unread"`
	IsMine           bool                    `json:"is_mine"`
	ReplyToMessageID *uint                   `json:"reply_to_message_id,omitempty"`
	ReplyToMessage   *ReplyMessageResponse   `json:"reply_to_message,omitempty"`
	// IsEdited  bool                    `json:"is_edited"`
}

type MessageReactionSummary struct {
	Emoji        string `json:"emoji"`
	Count        int    `json:"count"`
	ReactedByMe  bool   `json:"reacted_by_me"`
}

type ReplyMessageResponse struct {
	ID          uint                  `json:"id"`
	MessageType enums.ChatMessageType `json:"message_type"`
	Content     string                `json:"content"`
	MediaURL    *string               `json:"media_url,omitempty"`
	Sender      ChatParticipantResponse `json:"sender"`
	CreatedAt   time.Time             `json:"created_at"`
	IsEdited    bool                  `json:"is_edited"`
	IsDeleted   bool                  `json:"is_deleted"`
	DeletedAt   *time.Time            `json:"deleted_at,omitempty"`
}

type UserRoomsResponse struct {
	RoomID             uint                    `json:"room_id"`
	PetSitter          ChatParticipantResponse `json:"pet_sitter"`
	Status             enums.ChatRoomStatus    `json:"status"`               // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy          *enums.ChatBlockedBy    `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
	LastMessage        *string                 `json:"last_message,omitempty"`
	LastMessageType    enums.ChatMessageType   `json:"last_message_type"`
	LastMessageTime    *time.Time              `json:"last_message_time,omitempty"`
	UnreadMessageCount uint                    `json:"unread_message_count"`
}

type PetSitterRoomsResponse struct {
	RoomID             uint                    `json:"room_id"`
	User               ChatParticipantResponse `json:"user"`
	Status             enums.ChatRoomStatus    `json:"status"`               // PENDING | ACCEPTED | REJECTED | BLOCKED
	BlockedBy          *enums.ChatBlockedBy    `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
	LastMessage        *string                 `json:"last_message,omitempty"`
	LastMessageType    enums.ChatMessageType   `json:"last_message_type"`
	LastMessageTime    *time.Time              `json:"last_message_time,omitempty"`
	UnreadMessageCount int64                   `json:"unread_message_count"`
}

type SaveMessageResponse struct {
	ID        uint      `json:"id"`
	RoomID    uint      `json:"room_id"`	
	SenderID  uint      `json:"sender_id"`
	MessageType enums.ChatMessageType `json:"message_type"`
	Content   string    `json:"content"`
	MediaURL  *string   `json:"media_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	IsEdited  bool      `json:"is_edited"`
	IsDeleted bool      `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
