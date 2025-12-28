package service

import (
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainpostgres "hona/backend/internal/domain/ports/postgres"
	"log"
)

type ChatService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
}

func NewChatService(unitOfWork ports.UnitOfWork, userService usecase.UserService) *ChatService {
	return &ChatService{
		unitOfWork:  unitOfWork,
		userService: userService,
	}
}

func (cs *ChatService) CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	if info.PetSitterID == info.UserID {
		return chat.ChatRoomDetailsResponse{}, nil
	}
	log.Println("1")
	room, err := chatRepo.GetUserAndPetSitterRoom(info.UserID, info.PetSitterID)
	if err != nil {
		return chat.ChatRoomDetailsResponse{}, err
	}
	log.Println("2")
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

func (cs *ChatService) SaveMessage(info chat.SaveMessageRequest) (chat.SaveMessageResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	message := &entities.ChatMessage{
		RoomID:        info.RoomID,
		SenderID:      info.SenderID,
		Content:       info.Content,
		ReplyToMessageID: info.ReplyToMessageID,
	}
	err := chatRepo.CreateMessage(message)
	if err != nil {
		return chat.SaveMessageResponse{}, err	
	}
	return chat.SaveMessageResponse{
		ID:        message.ID,
		RoomID:    message.RoomID,
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}, nil
}

func (cs *ChatService) GetAllRooms(request chat.GetAllRoomsRequest) ([]chat.PetSitterRoomsResponse, int64, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	options := domainpostgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(request.Sort)
	rooms, totalCount, err := chatRepo.GetAllRooms(request.SenderID, options)
	if err != nil {
		return nil, 0, err
	}
	var roomDTOs []chat.PetSitterRoomsResponse
	for _, room := range rooms {
		unreadMessageCount, err := chatRepo.UnreadMessageCount(room.ID, request.SenderID, room.UserLastReadMessageID)
		if err != nil {
			return nil, 0, err
		}
		lastMessage, err := chatRepo.FindLastMessageByID(room.LastMessageID)
		if err != nil {
			return nil, 0, err
		}
		roomDTOs = append(roomDTOs, chat.PetSitterRoomsResponse{
			RoomID: room.ID,
			User: chat.ChatParticipantResponse{
				ID:        room.UserID,
				FirstName: room.User.FirstName,
				LastName:  room.User.LastName,
				// ProfilePic: room.User.ProfilePic,
				// IsOnline:   room.User.IsOnline
			},
			Status:             room.Status,
			BlockedBy:          room.BlockedBy,
			LastMessage:        &lastMessage.Content,
			LastMessageTime:    &lastMessage.CreatedAt,
			UnreadMessageCount: unreadMessageCount,
		})
	}
	return roomDTOs, totalCount, nil
}

func (cs *ChatService) GetRoomMessages(request *chat.GetRoomMessagesRequest) ([]chat.RoomMessagesResponse, int64, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	options := domainpostgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(request.Sort)
	room, err := chatRepo.GetRoomByID(request.RoomID)
	if err != nil {
		return nil, 0, err
	}
	if room == nil {
		// notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		return nil, 0, nil
	}
	messages, totalCount, err := chatRepo.GetMessagesByRoomID(request.RoomID, options)
	if err != nil {
		return nil, 0, err
	}
	var messageDTOs []chat.RoomMessagesResponse
	for _, message := range messages {
		sender, err := cs.userService.FindUserByID(message.SenderID)
		if err != nil {
			return nil, 0, err
		}
		messageDTOs = append(messageDTOs, chat.RoomMessagesResponse{
			ID:        message.ID,
			RoomID:    message.RoomID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
			Sender: chat.ChatParticipantResponse{
				FirstName:  sender.FirstName,
				LastName:   sender.LastName,
				// ProfilePic: sender.ProfilePic,
				// IsOnline:   sender.IsOnline,
			},
			IsUnread:         message.ID > *room.LastMessageID && room.UserID != message.SenderID,
			IsMine:           message.SenderID == room.UserID,
			ReplyToMessageID: message.ReplyToMessageID,
		})
	}
	return messageDTOs, totalCount, nil
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
	result := map[string]interface{}{
		"room_id":       room.ID,
		"user_id":       room.UserID,
		"pet_sitter_id": room.PetSitterID,
		"status":        room.Status.String(),
		"request_id":    room.RequestID,
	}
	return result, nil
}

func (cs *ChatService) AcceptRoom(petSitterID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return nil
	}
	if room.PetSitterID != petSitterID {
		return nil
	}
	room.Status = enums.C_Accepted
	return chatRepo.UpdateRoom(room)
}

func (cs *ChatService) RejectRoom(petSitterID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return nil
	}
	if room.PetSitterID != petSitterID {
		return nil
	}
	room.Status = enums.C_Rejected
	return chatRepo.UpdateRoom(room)
}

func (cs *ChatService) BlockRoom(request chat.BlockRoomRequest) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(request.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if room.Status == enums.C_Blocked {
		return nil
	}
	if request.SenderID != room.UserID && request.SenderID != room.PetSitterID {
		return exceptions.NewUnauthorizedError("You are not a participant of this chat room")
	}
	room.Status = enums.C_Blocked
	room.BlockedBy = &request.BlockedBy
	err = chatRepo.UpdateRoom(room)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ChatService) UnblockRoom(request chat.UnblockRoomRequest) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(request.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if request.SenderID != room.UserID && request.SenderID != room.PetSitterID {
		return exceptions.NewUnauthorizedError("You are not a participant of this chat room")
	}
	room.Status = enums.C_Accepted
	room.BlockedBy = nil
	err = chatRepo.UpdateRoom(room)
	if err != nil {
		return err
	}
	return nil
}

// func (cs *ChatService) GetRoomByID(id uint) (*entities.ChatRoom, error) {
// 	chatRepo := cs.unitOfWork.Factory().ChatRepository()
// 	room, err := chatRepo.FindRoomByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return room, nil
// }

func (cs *ChatService) MarkAsRead(request chat.MarkRoomReadRequest) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(request.RoomID)
	if err != nil {
		return err	
	}
	if room == nil {
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if room.Status != enums.C_Accepted {
		return nil
		// return exceptions.NewInvalidOperationError("Cannot mark messages as read in a non-accepted chat room")
	}
	isUser := request.SenderID == room.UserID
	isPetSitter := request.SenderID == room.PetSitterID
	if !isUser && !isPetSitter {
		return exceptions.NewUnauthorizedError("You are not a participant of this chat room")
	}
	message, err := chatRepo.FindLastMessageByID(&request.LastReadMessageID)
	if err != nil {
		return err
	}	
	if message == nil || message.RoomID != request.RoomID {
		return nil
	}
	if isUser {
		if room.UserLastReadMessageID !=nil && request.LastReadMessageID <= *room.UserLastReadMessageID {
			return nil
		}
		room.UserLastReadMessageID = &request.LastReadMessageID	
	} else if isPetSitter {
		if room.PetSitterLastReadMessageID != nil && request.LastReadMessageID <= *room.PetSitterLastReadMessageID {
			return nil
		}
		room.PetSitterLastReadMessageID = &request.LastReadMessageID
	}
	if err := chatRepo.UpdateRoom(room); err != nil {
		return err
	}
	return nil
}