package usecase

import "hona/backend/internal/application/dto/chat"

type ChatService interface {
	CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error)
	SaveMessage(info chat.SaveMessageRequest) (chat.SaveMessageResponse, error)
	GetAllRooms(request chat.GetAllRoomsRequest) ([]chat.PetSitterRoomsResponse, int64, error)
	BlockRoom(request chat.BlockRoomRequest) error
	UnblockRoom(request chat.UnblockRoomRequest) error
	GetRoomRequestInfo(roomID uint) (map[string]interface{}, error)
	AcceptRoom(petSitterID uint, roomID uint) error
	RejectRoom(petSitterID uint, roomID uint) error
	GetRoomMessages(request *chat.GetRoomMessagesRequest) ([]chat.RoomMessagesResponse, int64, error)
	// GetRoomByID(id uint) (*entities.ChatRoom, error)
}
