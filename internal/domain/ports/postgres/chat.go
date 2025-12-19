package domainpostgres

import "hona/backend/internal/domain/entities"

type ChatRepository interface {
	CreateRoom(room *entities.ChatRoom) error
	GetUserAndPetSitterRoom(userID uint, petSitterID uint) (*entities.ChatRoom, error)
	GetAllRoomsByUserID(userID uint) ([]*entities.ChatRoom, error)
	GetAllRoomsByPetSitterID(petSitterID uint) ([]*entities.ChatRoom, error)
	UpdateRoomStatus(roomID uint, status uint) error
	UpdateRoomBlockedBy(roomID uint, blockedBy string) error
	GetRoomByID(roomID uint) (*entities.ChatRoom, error)
	GetRequestIDByRoomID(roomID uint) (uint, error)
	SaveMessage(message *entities.ChatMessage) error
}
