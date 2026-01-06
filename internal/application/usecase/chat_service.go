package usecase

import "hona/backend/internal/application/dto/chat"

type ChatService interface {
	CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error)
	SaveMessage(info chat.SaveMessageRequest) (chat.SaveMessageResponse, error)
	GetAllRooms(request chat.GetAllRoomsRequest) ([]chat.PetSitterRoomsResponse, int64, error)
	GetUserRooms(request chat.GetAllRoomsRequest) ([]chat.UserRoomsResponse, int64, error)
	EditMessage(request chat.EditMessageRequest) (chat.SaveMessageResponse, error)
	DeleteMessage(request chat.DeleteMessageRequest) (chat.SaveMessageResponse, error)
	AddReaction(request chat.ReactionRequest) ([]chat.MessageReactionSummary, error)
	RemoveReaction(request chat.ReactionRequest) ([]chat.MessageReactionSummary, error)
	BlockRoom(request chat.BlockRoomRequest) error
	UnblockRoom(request chat.UnblockRoomRequest) error
	GetRoomRequestInfo(roomID uint, senderID uint) (map[string]interface{}, error)
	AcceptRoom(petSitterID uint, roomID uint) error
	RejectRoom(petSitterID uint, roomID uint) error
	GetRoomMessages(request *chat.GetRoomMessagesRequest) ([]chat.RoomMessagesResponse, int64, error)
	MarkAsRead(request chat.MarkRoomReadRequest) error
	// GetRoomByID(id uint) (*entities.ChatRoom, error)
}
