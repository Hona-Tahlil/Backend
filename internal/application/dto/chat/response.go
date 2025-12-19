package chat

import "time"

type ChatRoomDetailsResponse struct {
	RoomID    uint                    `json:"room_id"`
	User      ChatParticipantResponse `json:"user"`
	Petsitter ChatParticipantResponse `json:"petsitter"`
	Status    string                  `json:"status"` // PENDING | ACCEPTED | REJECTED | BLOCKED
	// IsAccepted bool                    `json:"is_accepted"`
	BlockedBy *string `json:"blocked_by,omitempty"` // "USER" | "PETSITTER"
}

type ChatParticipantResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	IsOnline  bool   `json:"is_online"`
}

type RoomMessagesResponse struct {
	ID        uint                    `json:"id"`
	Sender    ChatParticipantResponse `json:"sender"`
	Content   string                  `json:"content"`
	TimeStamp time.Time               `json:"timeStamp"`
}
