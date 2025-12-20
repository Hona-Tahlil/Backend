package websocket

import (
	"encoding/json"
	"time"
)

const (
	MessageTypeChat     = "CHAT"
	MessageTypePresence = "PRESENCE"
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

type PresencePayload struct {
	UserID    uint      `json:"user_id"`
	IsOnline  bool      `json:"is_online"`
	Timestamp time.Time `json:"timestamp"`
}
