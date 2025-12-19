package usecase

import "hona/backend/internal/application/dto/chat"

type ChatService interface {
	// Define the methods that the ChatService should have
	CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error)
	SaveMessage(info chat.SaveMessageRequest) error
	GetAllRooms(userID uint) ([]chat.ChatRoomDetailsResponse, error)
	BlockRoom(userID uint, roomID uint) error
	UnblockRoom(userID uint, roomID uint) error
	GetRoomRequestInfo(roomID uint) (map[string]interface{}, error)
	AcceptRoom(petSitterID uint, userID uint) error
	RejectRoom(petSitterID uint, roomID uint) error
	// GetRoomByID(id uint) (*entities.ChatRoom, error)
}
