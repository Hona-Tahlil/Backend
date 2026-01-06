package domainpostgres

import (
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
)

type ChatRepository interface {
	CreateRoom(room *entities.ChatRoom) error
	GetUserAndPetSitterRoom(userID uint, petSitterID uint) (*entities.ChatRoom, error)
	GetAllRoomsByUserID(userID uint) ([]*entities.ChatRoom, error)
	GetAllRoomsByPetSitterID(petSitterID uint) ([]*entities.ChatRoom, error)
	UpdateRoomStatus(roomID uint, status uint) error
	UpdateRoomBlockedBy(roomID uint, blockedBy string) error
	GetRoomByID(roomID uint) (*entities.ChatRoom, error)
	GetRequestIDByRoomID(roomID uint) (*uint, error)
	CreateMessage(message *entities.ChatMessage) error
	UpdateRoom(room *entities.ChatRoom) error
	GetAllRooms(senderID uint, options *postgres.QueryOptions) ([]*entities.ChatRoom, int64, error)
	UnreadMessageCount(roomID uint, senderID uint, lastReadMessageID *uint) (int64, error)
	FindLastMessageByID(messageID *uint) (*entities.ChatMessage, error)
	FindLastActiveMessageInRoom(roomID uint) (*entities.ChatMessage, error)
	FindMessageByID(messageID uint) (*entities.ChatMessage, error)
	FindReaction(messageID uint, userID uint, emoji string) (*entities.ChatMessageReaction, error)
	AddReaction(reaction *entities.ChatMessageReaction) error
	RemoveReaction(messageID uint, userID uint, emoji string) error
	GetReactionsByMessageIDs(messageIDs []uint) (map[uint][]*entities.ChatMessageReaction, error)
	GetMessagesByRoomID(roomID uint, options *postgres.QueryOptions) ([]*entities.ChatMessage, int64, error)
	UpdateMessage(message *entities.ChatMessage) error
	DeleteMessageByID(messageID uint) error
}
