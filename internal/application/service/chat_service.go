package service

import (
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/ports"
)

type ChatService struct {
	unitOfWork ports.UnitOfWork
}

func NewChatService(unitOfWork ports.UnitOfWork) *ChatService {
	return &ChatService{
		unitOfWork: unitOfWork,
	}
}

func (cs *ChatService) CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	if info.PetSitterID == info.UserID {
		return chat.ChatRoomDetailsResponse{}, nil
	}
	room, err := chatRepo.GetUserAndPetSitterRoom(info.UserID, info.PetSitterID)
	if err != nil {
		return chat.ChatRoomDetailsResponse{}, err
	}
	if room == nil {
		newRoom := &entities.ChatRoom{
			UserID:      info.UserID,
			PetSitterID: info.PetSitterID,
		}
		err := chatRepo.CreateRoom(newRoom)
		if err != nil {
			return chat.ChatRoomDetailsResponse{}, err
		}
	}
	return chat.ChatRoomDetailsResponse{RoomID: room.ID}, nil
}

func (cs *ChatService) SaveMessage(info chat.SaveMessageRequest) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	message := &entities.ChatMessage{
		RoomID:   info.RoomID,
		SenderID: info.UserID,
		Content:  info.Content,
	}
	return chatRepo.SaveMessage(message)
}

func (cs *ChatService) GetAllRooms(userID uint) ([]chat.ChatRoomDetailsResponse, error) {
	// This can be called by either user or petsitter controller
	// For now return empty slice; actual implementation would fetch from DB and map to DTO
	return []chat.ChatRoomDetailsResponse{}, nil
}

func (cs *ChatService) BlockRoom(userID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	// Update room's blocked_by field to indicate who blocked it
	return chatRepo.UpdateRoomBlockedBy(roomID, "USER")
}

func (cs *ChatService) UnblockRoom(userID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	// Clear the blocked_by field
	return chatRepo.UpdateRoomBlockedBy(roomID, "")
}

func (cs *ChatService) GetRoomRequestInfo(roomID uint) (map[string]interface{}, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return map[string]interface{}{}, nil
	}
	// Return room details as a map
	result := map[string]interface{}{
		"room_id":       room.ID,
		"user_id":       room.UserID,
		"pet_sitter_id": room.PetSitterID,
		"status":        room.Status.String(),
		"request_id":    room.RequestID,
	}
	return result, nil
}

func (cs *ChatService) AcceptRoom(petSitterID uint, userID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetUserAndPetSitterRoom(userID, petSitterID)
	if err != nil {
		return err
	}
	if room == nil {
		return nil
	}
	return chatRepo.UpdateRoomStatus(room.ID, uint(enums.C_Accepted))
}

func (cs *ChatService) RejectRoom(petSitterID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	return chatRepo.UpdateRoomStatus(roomID, uint(enums.C_Rejected))
}

// func (cs *ChatService) GetRoomByID(id uint) (*entities.ChatRoom, error) {
// 	chatRepo := cs.unitOfWork.Factory().ChatRepository()
// 	room, err := chatRepo.FindRoomByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return room, nil
// }
