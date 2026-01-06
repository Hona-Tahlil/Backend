package websocket

import (
	"encoding/json"
	"time"
)

const (
	MessageTypeChat     = "CHAT"
	MessageTypePresence = "PRESENCE"
	MessageTypeRead     = "READ"
	MessageTypeEdit     = "EDIT"
	MessageTypeDelete   = "DELETE"
	MessageTypeReaction = "REACTION"
)

type Message struct {
	Type      string          `json:"type"`
	RoomID    uint            `json:"room_id"`
	SenderID  uint            `json:"sender_id,omitempty"`
	MessageID uint            `json:"message_id,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Content   json.RawMessage `json:"content,omitempty"`

	Client *Client `json:"-"`
}

type ChatPayload struct {
	MessageType    string `json:"message_type"`
	Text           string `json:"text,omitempty"`
	ImageBase64    string `json:"image_base64,omitempty"`
	ImageMime      string `json:"image_mime,omitempty"`
	ImageURL       string `json:"image_url,omitempty"`
	ReplyToMessageID *uint `json:"reply_to_message_id,omitempty"`
	IsUnread       bool   `json:"is_unread"`
	IsMine         bool   `json:"is_mine"`
}

type PresencePayload struct {
	UserID    uint       `json:"user_id"`
	IsOnline  bool       `json:"is_online"`
	LastSeen  *time.Time `json:"last_seen,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

type ReadPayload struct {
	ReaderID          uint      `json:"reader_id"`
	LastReadMessageID uint      `json:"last_read_message_id"`
	Timestamp         time.Time `json:"timestamp"`
}

type EditPayload struct {
	MessageID   uint   `json:"message_id"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
	MediaURL    *string `json:"media_url,omitempty"`
	IsEdited    bool   `json:"is_edited"`
	EditedAt    time.Time `json:"edited_at"`
}

type DeletePayload struct {
	MessageID uint       `json:"message_id"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ReactionPayload struct {
	MessageID uint      `json:"message_id"`
	Emoji     string    `json:"emoji"`
	Action    string    `json:"action"`
	UserID    uint      `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
