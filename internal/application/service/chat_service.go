package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"hona/backend/internal/application/dto/chat"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainstorage "hona/backend/internal/domain/storage"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
	"hona/backend/internal/infrastructure/websocket"
	"sort"
	"strings"
	"time"

	"log"
)

type ChatService struct {
	unitOfWork  ports.UnitOfWork
	userService usecase.UserService
	storage     domainstorage.Storage
	hub         *websocket.Hub
	requestService usecase.RequestService
}

func NewChatService(unitOfWork ports.UnitOfWork, userService usecase.UserService, storage domainstorage.Storage, hub *websocket.Hub, requestService usecase.RequestService) *ChatService {
	return &ChatService{
		unitOfWork:  unitOfWork,
		userService: userService,
		storage:     storage,
		hub:         hub,
		requestService: requestService,
	}
}

func (cs *ChatService) CreateOrGetRoom(info chat.CreateOrGetUserRoomRequest) (chat.ChatRoomDetailsResponse, error) {
	log.Println("Creating or getting room")
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	if info.PetSitterID == info.UserID {
		return chat.ChatRoomDetailsResponse{}, nil
	}
	user, err := cs.userService.FindUserByID(info.UserID)
	if err != nil {
		return chat.ChatRoomDetailsResponse{}, err
	}
	// petSitter, err := cs.userService.FindUserByID(info.PetSitterID)
	// if err != nil {
	// 	return chat.ChatRoomDetailsResponse{}, err
	// }
	// petSitterRepo := cs.unitOfWork.Factory().PetSitterRepository()
	// petSitter, err := petSitterRepo.FindPetSitterByID(info.PetSitterID)
	// if err != nil {
	// 	return chat.ChatRoomDetailsResponse{}, err
	// }
	// if petSitter == nil {
	// 	return chat.ChatRoomDetailsResponse{}, exceptions.NewNotFoundError("PetSitter")
	// }

	petSitterUser, err := cs.userService.FindUserByID(info.PetSitterID)
	if err != nil {
		return chat.ChatRoomDetailsResponse{}, err
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
			Status:      enums.C_Pending,
		}
		err := chatRepo.CreateRoom(newRoom)
		if err != nil {
			return chat.ChatRoomDetailsResponse{}, err
		}
		room = newRoom

	}

	var blockedBy *string
	if room.BlockedBy != nil {
		blockedByValue := room.BlockedBy.String()
		blockedBy = &blockedByValue
	}

	userProfilePic := ""
	if user.PictureLink != nil {
		userProfilePic = *user.PictureLink
	}
	petSitterProfilePic := ""
	if petSitterUser.PictureLink != nil {
		petSitterProfilePic = *petSitterUser.PictureLink
	}

	return chat.ChatRoomDetailsResponse{
		RoomID: room.ID,
		User: chat.ChatParticipantResponse{
			ID:         user.ID,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			ProfilePic: userProfilePic,
			IsOnline:   cs.isUserOnline(user.ID),
			LastSeen:   cs.lastSeen(user.ID),
		},
		PetSitter: chat.ChatParticipantResponse{
			ID:         petSitterUser.ID,
			FirstName:  petSitterUser.FirstName,
			LastName:   petSitterUser.LastName,
			ProfilePic: petSitterProfilePic,
			IsOnline:   cs.isUserOnline(petSitterUser.ID),
			LastSeen:   cs.lastSeen(petSitterUser.ID),
		},
		Status:    room.Status.String(),
		BlockedBy: blockedBy,
	}, nil
}

func (cs *ChatService) SaveMessage(info chat.SaveMessageRequest) (chat.SaveMessageResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := cs.getRoomForParticipant(info.RoomID, info.SenderID)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if err := ensureRoomAccepted(room); err != nil {
		return chat.SaveMessageResponse{}, err
	}
	messageType := info.MessageType
	if messageType == "" {
		messageType = enums.ChatMessageText
	}
	if info.ReplyToMessageID != nil {
		ref, err := chatRepo.FindMessageByID(*info.ReplyToMessageID)
		if err != nil {
			return chat.SaveMessageResponse{}, err
		}
		if ref == nil || ref.RoomID != info.RoomID {
			return chat.SaveMessageResponse{}, errors.New("invalid reply_to_message_id")
		}
	}

	var mediaKey *string
	var mediaURL *string
	if messageType == enums.ChatMessageImage {
		contentType, data, err := decodeBase64Image(info.MediaBase64, info.MediaMime)
		if err != nil {
			return chat.SaveMessageResponse{}, err
		}
		if len(data) > 2*1024*1024 {
			return chat.SaveMessageResponse{}, errors.New("image size exceeds 2MB")
		}
		key, err := buildChatMediaKey(info.RoomID, contentType)
		if err != nil {
			return chat.SaveMessageResponse{}, err
		}
		if err := cs.storage.UploadBytes(enums.ChatMedia, key, contentType, data); err != nil {
			return chat.SaveMessageResponse{}, err
		}
		mediaKey = &key
		url, err := cs.storage.GetPresignedURL(enums.ChatMedia, key, time.Minute*15)
		if err == nil {
			mediaURL = &url
		}
	}
	message := &entities.ChatMessage{
		RoomID:           info.RoomID,
		SenderID:         info.SenderID,
		Content:          info.Content,
		MessageType:      messageType,
		MediaKey:         mediaKey,
		ReplyToMessageID: info.ReplyToMessageID,
	}
	err = chatRepo.CreateMessage(message)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if room, err := chatRepo.GetRoomByID(info.RoomID); err == nil && room != nil {
		room.LastMessageID = &message.ID
		_ = chatRepo.UpdateRoom(room)
	}
	return chat.SaveMessageResponse{
		ID:          message.ID,
		RoomID:      message.RoomID,
		SenderID:    message.SenderID,
		MessageType: message.MessageType,
		Content:     message.Content,
		MediaURL:    mediaURL,
		CreatedAt:   message.CreatedAt,
		IsEdited:    message.IsEdited,
		IsDeleted:   message.IsDeleted,
		DeletedAt:   message.DeletedAtTime,
	}, nil
}

func (cs *ChatService) GetAllRooms(request chat.GetAllRoomsRequest) ([]chat.PetSitterRoomsResponse, int64, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)
	rooms, totalCount, err := chatRepo.GetAllRooms(request.SenderID, options)
	if err != nil {
		return nil, 0, err
	}
	var roomDTOs []chat.PetSitterRoomsResponse
	for _, room := range rooms {
		lastReadMessageID := lastReadMessageIDForSender(room, request.SenderID)
		unreadMessageCount, err := chatRepo.UnreadMessageCount(room.ID, request.SenderID, lastReadMessageID)
		if err != nil {
			return nil, 0, err
		}
		lastMessage, err := chatRepo.FindLastMessageByID(room.LastMessageID)
		if err != nil {
			return nil, 0, err
		}
		user, err := cs.userService.FindUserByID(room.UserID)
		if err != nil {
			return nil, 0, err
		}
		userProfilePic := ""
		if user.PictureLink != nil {
			userProfilePic = *user.PictureLink
		}
		lastMessageType := enums.ChatMessageText
		var lastMessageContent *string
		var lastMessageTime *time.Time
		if lastMessage != nil {
			lastMessageType = lastMessage.MessageType
			if lastMessage.MessageType == enums.ChatMessageImage {
				if strings.TrimSpace(lastMessage.Content) != "" {
					lastMessageContent = &lastMessage.Content
				} else {
					placeholder := "[image]"
					lastMessageContent = &placeholder
				}
			} else {
				lastMessageContent = &lastMessage.Content
			}
			lastMessageTime = &lastMessage.CreatedAt
		}
		roomDTOs = append(roomDTOs, chat.PetSitterRoomsResponse{
			RoomID: room.ID,
			User: chat.ChatParticipantResponse{
				ID:         user.ID,
				FirstName:  user.FirstName,
				LastName:   user.LastName,
				ProfilePic: userProfilePic,
				IsOnline:   cs.isUserOnline(user.ID),
				LastSeen:   cs.lastSeen(user.ID),
			},
			Status:             room.Status,
			BlockedBy:          room.BlockedBy,
			LastMessage:        lastMessageContent,
			LastMessageType:    lastMessageType,
			LastMessageTime:    lastMessageTime,
			UnreadMessageCount: unreadMessageCount,
		})
	}
	return roomDTOs, totalCount, nil
}

func (cs *ChatService) GetUserRooms(request chat.GetAllRoomsRequest) ([]chat.UserRoomsResponse, int64, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)
	log.Println("Fetching user rooms for sender:", request.SenderID)
	rooms, totalCount, err := chatRepo.GetAllRooms(request.SenderID, options)
	log.Println("Total rooms found:", totalCount)
	if err != nil {
		return nil, 0, err
	}
	var roomDTOs []chat.UserRoomsResponse
	for _, room := range rooms {
		lastReadMessageID := lastReadMessageIDForSender(room, request.SenderID)
		unreadMessageCount, err := chatRepo.UnreadMessageCount(room.ID, request.SenderID, lastReadMessageID)
		if err != nil {
			return nil, 0, err
		}
		lastMessage, err := chatRepo.FindLastMessageByID(room.LastMessageID)
		if err != nil {
			return nil, 0, err
		}
		petSitterUser, err := cs.userService.FindUserByID(room.PetSitterID)
		if err != nil {
			return nil, 0, err
		}
		petSitterProfilePic := ""
		if petSitterUser.PictureLink != nil {
			petSitterProfilePic = *petSitterUser.PictureLink
		}
		lastMessageType := enums.ChatMessageText
		var lastMessageContent *string
		var lastMessageTime *time.Time
		if lastMessage != nil {
			lastMessageType = lastMessage.MessageType
			if lastMessage.MessageType == enums.ChatMessageImage {
				if strings.TrimSpace(lastMessage.Content) != "" {
					lastMessageContent = &lastMessage.Content
				} else {
					placeholder := "[image]"
					lastMessageContent = &placeholder
				}
			} else {
				lastMessageContent = &lastMessage.Content
			}
			lastMessageTime = &lastMessage.CreatedAt
		}
		roomDTOs = append(roomDTOs, chat.UserRoomsResponse{
			RoomID: room.ID,
			PetSitter: chat.ChatParticipantResponse{
				ID:         petSitterUser.ID,
				FirstName:  petSitterUser.FirstName,
				LastName:   petSitterUser.LastName,
				ProfilePic: petSitterProfilePic,
				IsOnline:   cs.isUserOnline(petSitterUser.ID),
				LastSeen:   cs.lastSeen(petSitterUser.ID),
			},
			Status:             room.Status,
			BlockedBy:          room.BlockedBy,
			LastMessage:        lastMessageContent,
			LastMessageType:    lastMessageType,
			LastMessageTime:    lastMessageTime,
			UnreadMessageCount: uint(unreadMessageCount),
		})
	}
	return roomDTOs, totalCount, nil
}

func (cs *ChatService) GetRoomMessages(request *chat.GetRoomMessagesRequest) ([]chat.RoomMessagesResponse, int64, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)
	room, err := cs.getRoomForParticipant(request.RoomID, request.SenderID)
	if err != nil {
		return nil, 0, err
	}
	messages, totalCount, err := chatRepo.GetMessagesByRoomID(request.RoomID, options)
	if err != nil {
		return nil, 0, err
	}
	messageIDs := make([]uint, 0, len(messages))
	for _, message := range messages {
		messageIDs = append(messageIDs, message.ID)
	}
	reactionsByMessageID, err := chatRepo.GetReactionsByMessageIDs(messageIDs)
	if err != nil {
		return nil, 0, err
	}
	var messageDTOs []chat.RoomMessagesResponse
	for _, message := range messages {
		sender, err := cs.userService.FindUserByID(message.SenderID)
		if err != nil {
			return nil, 0, err
		}
		lastReadMessageID := lastReadMessageIDForSender(room, request.SenderID)
		isUnread := false
		if message.SenderID != request.SenderID {
			if lastReadMessageID == nil {
				isUnread = true
			} else {
				isUnread = message.ID > *lastReadMessageID
			}
		}
		var mediaURL *string
		content := message.Content
		if message.IsDeleted {
			content = ""
		} else if message.MessageType == enums.ChatMessageImage && message.MediaKey != nil {
			url, err := cs.storage.GetPresignedURL(enums.ChatMedia, *message.MediaKey, time.Minute*15)
			if err == nil {
				mediaURL = &url
			}
		}
		var reactions []chat.MessageReactionSummary
		if !message.IsDeleted {
			reactions = buildReactionSummary(reactionsByMessageID[message.ID], request.SenderID)
		}
		var replyToMessage *chat.ReplyMessageResponse
		if message.ReplyToMessageID != nil {
			reply, err := chatRepo.FindMessageByID(*message.ReplyToMessageID)
			if err != nil {
				return nil, 0, err
			}
			if reply != nil && reply.RoomID == message.RoomID {
				replySender, err := cs.userService.FindUserByID(reply.SenderID)
				if err != nil {
					return nil, 0, err
				}
				replyProfilePic := ""
				if replySender.PictureLink != nil {
					replyProfilePic = *replySender.PictureLink
				}
				var replyMediaURL *string
				replyContent := reply.Content
				if reply.IsDeleted {
					replyContent = ""
				} else if reply.MessageType == enums.ChatMessageImage && reply.MediaKey != nil {
					url, err := cs.storage.GetPresignedURL(enums.ChatMedia, *reply.MediaKey, time.Minute*15)
					if err == nil {
						replyMediaURL = &url
					}
				}
				replyToMessage = &chat.ReplyMessageResponse{
					ID:          reply.ID,
					MessageType: reply.MessageType,
					Content:     replyContent,
					MediaURL:    replyMediaURL,
					Sender: chat.ChatParticipantResponse{
						ID:         replySender.ID,
						FirstName:  replySender.FirstName,
						LastName:   replySender.LastName,
						ProfilePic: replyProfilePic,
						IsOnline:   cs.isUserOnline(replySender.ID),
						LastSeen:   cs.lastSeen(replySender.ID),
					},
					CreatedAt: reply.CreatedAt,
					IsEdited:  reply.IsEdited,
					IsDeleted: reply.IsDeleted,
					DeletedAt: reply.DeletedAtTime,
				}
			}
		}
		messageDTOs = append(messageDTOs, chat.RoomMessagesResponse{
			ID:        message.ID,
			RoomID:    message.RoomID,
			MessageType: message.MessageType,
			Content:   content,
			MediaURL:  mediaURL,
			CreatedAt: message.CreatedAt,
			IsEdited:  message.IsEdited,
			IsDeleted: message.IsDeleted,
			DeletedAt: message.DeletedAtTime,
			Reactions: reactions,
			Sender: chat.ChatParticipantResponse{
				FirstName: sender.FirstName,
				LastName:  sender.LastName,
				// ProfilePic: sender.ProfilePic,
				IsOnline: cs.isUserOnline(sender.ID),
				LastSeen: cs.lastSeen(sender.ID),
			},
			IsUnread:         isUnread,
			IsMine:           message.SenderID == request.SenderID,
			ReplyToMessageID: message.ReplyToMessageID,
			ReplyToMessage:   replyToMessage,
		})
	}
	return messageDTOs, totalCount, nil
}

func (cs *ChatService) AddReaction(request chat.ReactionRequest) ([]chat.MessageReactionSummary, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	if request.Emoji == "" {
		return nil, errors.New("emoji is required")
	}
	room, err := cs.getRoomForParticipant(request.RoomID, request.SenderID)
	if err != nil {
		return nil, err
	}
	if err := ensureRoomAccepted(room); err != nil {
		return nil, err
	}
	message, err := chatRepo.FindMessageByID(request.MessageID)
	if err != nil {
		return nil, err
	}
	if message == nil || message.RoomID != request.RoomID {
		return nil, exceptions.NewNotFoundError("ChatMessage")
	}
	if message.IsDeleted {
		return nil, errors.New("message is deleted")
	}
	existing, err := chatRepo.FindReaction(request.MessageID, request.SenderID, request.Emoji)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		reaction := &entities.ChatMessageReaction{
			MessageID: request.MessageID,
			UserID:    request.SenderID,
			Emoji:     request.Emoji,
		}
		if err := chatRepo.AddReaction(reaction); err != nil {
			return nil, err
		}
	}
	reactionsByMessageID, err := chatRepo.GetReactionsByMessageIDs([]uint{request.MessageID})
	if err != nil {
		return nil, err
	}
	return buildReactionSummary(reactionsByMessageID[request.MessageID], request.SenderID), nil
}

func (cs *ChatService) RemoveReaction(request chat.ReactionRequest) ([]chat.MessageReactionSummary, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	if request.Emoji == "" {
		return nil, errors.New("emoji is required")
	}
	room, err := cs.getRoomForParticipant(request.RoomID, request.SenderID)
	if err != nil {
		return nil, err
	}
	if err := ensureRoomAccepted(room); err != nil {
		return nil, err
	}
	message, err := chatRepo.FindMessageByID(request.MessageID)
	if err != nil {
		return nil, err
	}
	if message == nil || message.RoomID != request.RoomID {
		return nil, exceptions.NewNotFoundError("ChatMessage")
	}
	if message.IsDeleted {
		return nil, errors.New("message is deleted")
	}
	if err := chatRepo.RemoveReaction(request.MessageID, request.SenderID, request.Emoji); err != nil {
		return nil, err
	}
	reactionsByMessageID, err := chatRepo.GetReactionsByMessageIDs([]uint{request.MessageID})
	if err != nil {
		return nil, err
	}
	return buildReactionSummary(reactionsByMessageID[request.MessageID], request.SenderID), nil
}

func (cs *ChatService) EditMessage(request chat.EditMessageRequest) (chat.SaveMessageResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := cs.getRoomForParticipant(request.RoomID, request.SenderID)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if err := ensureRoomAccepted(room); err != nil {
		return chat.SaveMessageResponse{}, err
	}
	message, err := chatRepo.FindMessageByID(request.MessageID)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if message == nil || message.RoomID != request.RoomID {
		return chat.SaveMessageResponse{}, exceptions.NewNotFoundError("ChatMessage")
	}
	if message.SenderID != request.SenderID {
		return chat.SaveMessageResponse{}, exceptions.NewUnauthorizedError("You are not the sender of this message")
	}
	if message.IsDeleted {
		return chat.SaveMessageResponse{}, errors.New("message is deleted")
	}

	messageType := message.MessageType
	if request.MessageType != "" {
		messageType = request.MessageType
	}

	var mediaURL *string
	if messageType == enums.ChatMessageImage {
		if request.MediaBase64 != "" {
			contentType, data, err := decodeBase64Image(request.MediaBase64, request.MediaMime)
			if err != nil {
				return chat.SaveMessageResponse{}, err
			}
			if len(data) > 2*1024*1024 {
				return chat.SaveMessageResponse{}, errors.New("image size exceeds 2MB")
			}
			key, err := buildChatMediaKey(request.RoomID, contentType)
			if err != nil {
				return chat.SaveMessageResponse{}, err
			}
			if err := cs.storage.UploadBytes(enums.ChatMedia, key, contentType, data); err != nil {
				return chat.SaveMessageResponse{}, err
			}
			message.MediaKey = &key
		}
		if message.MediaKey != nil {
			url, err := cs.storage.GetPresignedURL(enums.ChatMedia, *message.MediaKey, time.Minute*15)
			if err == nil {
				mediaURL = &url
			}
		}
	} else {
		message.MediaKey = nil
	}

	message.MessageType = messageType
	message.Content = request.Content
	now := time.Now().UTC()
	message.IsEdited = true
	message.EditedAt = &now
	if err := chatRepo.UpdateMessage(message); err != nil {
		return chat.SaveMessageResponse{}, err
	}
	return chat.SaveMessageResponse{
		ID:          message.ID,
		RoomID:      message.RoomID,
		SenderID:    message.SenderID,
		MessageType: message.MessageType,
		Content:     message.Content,
		MediaURL:    mediaURL,
		CreatedAt:   message.CreatedAt,
		IsEdited:    message.IsEdited,
		IsDeleted:   message.IsDeleted,
		DeletedAt:   message.DeletedAtTime,
	}, nil
}

func (cs *ChatService) DeleteMessage(request chat.DeleteMessageRequest) (chat.SaveMessageResponse, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := cs.getRoomForParticipant(request.RoomID, request.SenderID)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if err := ensureRoomAccepted(room); err != nil {
		return chat.SaveMessageResponse{}, err
	}
	message, err := chatRepo.FindMessageByID(request.MessageID)
	if err != nil {
		return chat.SaveMessageResponse{}, err
	}
	if message == nil || message.RoomID != request.RoomID {
		return chat.SaveMessageResponse{}, exceptions.NewNotFoundError("ChatMessage")
	}
	if message.SenderID != request.SenderID {
		return chat.SaveMessageResponse{}, exceptions.NewUnauthorizedError("You are not the sender of this message")
	}
	if !message.IsDeleted {
		now := time.Now().UTC()
		message.IsDeleted = true
		message.DeletedAtTime = &now
		message.Content = ""
		message.MediaKey = nil
		if err := chatRepo.UpdateMessage(message); err != nil {
			return chat.SaveMessageResponse{}, err
		}
	}
	if room, err := chatRepo.GetRoomByID(request.RoomID); err == nil && room != nil && room.LastMessageID != nil && *room.LastMessageID == message.ID {
		lastActive, err := chatRepo.FindLastActiveMessageInRoom(request.RoomID)
		if err == nil {
			if lastActive != nil {
				room.LastMessageID = &lastActive.ID
			} else {
				room.LastMessageID = nil
			}
			_ = chatRepo.UpdateRoom(room)
		}
	}
	return chat.SaveMessageResponse{
		ID:          message.ID,
		RoomID:      message.RoomID,
		SenderID:    message.SenderID,
		MessageType: message.MessageType,
		Content:     message.Content,
		MediaURL:    nil,
		CreatedAt:   message.CreatedAt,
		IsEdited:    message.IsEdited,
		IsDeleted:   message.IsDeleted,
		DeletedAt:   message.DeletedAtTime,
	}, nil
}

func lastReadMessageIDForSender(room *entities.ChatRoom, senderID uint) *uint {
	if room == nil {
		return nil
	}
	if senderID == room.UserID {
		return room.UserLastReadMessageID
	}
	if senderID == room.PetSitterID {
		return room.PetSitterLastReadMessageID
	}
	return nil
}

func (cs *ChatService) isUserOnline(userID uint) bool {
	if cs.hub == nil {
		return false
	}
	return cs.hub.IsUserOnline(userID)
}

func (cs *ChatService) lastSeen(userID uint) *time.Time {
	if cs.hub == nil {
		return nil
	}
	if cs.hub.IsUserOnline(userID) {
		return nil
	}
	return cs.hub.LastSeen(userID)
}

func buildReactionSummary(reactions []*entities.ChatMessageReaction, senderID uint) []chat.MessageReactionSummary {
	if len(reactions) == 0 {
		return nil
	}
	byEmoji := make(map[string]*chat.MessageReactionSummary)
	for _, reaction := range reactions {
		summary, ok := byEmoji[reaction.Emoji]
		if !ok {
			summary = &chat.MessageReactionSummary{
				Emoji: reaction.Emoji,
			}
			byEmoji[reaction.Emoji] = summary
		}
		summary.Count++
		if reaction.UserID == senderID {
			summary.ReactedByMe = true
		}
	}
	emojis := make([]string, 0, len(byEmoji))
	for emoji := range byEmoji {
		emojis = append(emojis, emoji)
	}
	sort.Strings(emojis)
	out := make([]chat.MessageReactionSummary, 0, len(emojis))
	for _, emoji := range emojis {
		out = append(out, *byEmoji[emoji])
	}
	return out
}

func decodeBase64Image(raw string, mime string) (string, []byte, error) {
	if raw == "" {
		return "", nil, errors.New("image payload is empty")
	}
	data := raw
	contentType := mime
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return "", nil, errors.New("invalid data url")
		}
		meta := parts[0]
		data = parts[1]
		if contentType == "" {
			contentType = strings.TrimPrefix(strings.SplitN(meta, ";", 2)[0], "data:")
		}
	}
	if contentType == "" {
		return "", nil, errors.New("image mime type is required")
	}
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return "", nil, errors.New("unsupported image mime type")
	}
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", nil, errors.New("invalid base64 image data")
	}
	return contentType, decoded, nil
}

func buildChatMediaKey(roomID uint, mime string) (string, error) {
	ext := ""
	switch mime {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/webp":
		ext = "webp"
	default:
		return "", errors.New("unsupported image mime type")
	}
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	nonce := fmt.Sprintf("%x", randomBytes)
	return fmt.Sprintf("chat/%d/%d_%s.%s", roomID, time.Now().UnixNano(), nonce, ext), nil
}

func (cs *ChatService) GetRoomRequestInfo(roomID uint, senderID uint) (map[string]interface{}, error) {
	room, err := cs.getRoomForParticipant(roomID, senderID)
	if err != nil {
		return nil, err
	}
	requests, err := cs.requestService.GetRequestsBetween(request.GetRequestsBetweenRequest{
		UserID:          room.UserID,
		PetSitterUserID: room.PetSitterID,
	})
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"room_id":       room.ID,
		"user_id":       room.UserID,
		"pet_sitter_id": room.PetSitterID,
		"status":        room.Status.String(),
		"request_id":    room.RequestID,
		"requests":      requests,
	}
	return result, nil
}

func (cs *ChatService) getRoomForParticipant(roomID uint, senderID uint) (*entities.ChatRoom, error) {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, exceptions.NewNotFoundError("ChatRoom")
	}
	if senderID != room.UserID && senderID != room.PetSitterID {
		return nil, exceptions.NewUnauthorizedError("You are not a participant of this chat room")
	}
	return room, nil
}

func ensureRoomAccepted(room *entities.ChatRoom) error {
	if room == nil {
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if room.Status != enums.C_Accepted {
		return exceptions.NewAccessDeniedError("Chat room is not accepted")
	}
	return nil
}

func (cs *ChatService) AcceptRoom(petSitterID uint, roomID uint) error {
	chatRepo := cs.unitOfWork.Factory().ChatRepository()
	room, err := chatRepo.GetRoomByID(roomID)
	
	if err != nil {
		return err
	}
	if room == nil {
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if room.PetSitterID != petSitterID {
		return exceptions.NewUnauthorizedError("You are not a participant of this chat room")
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
		return exceptions.NewNotFoundError("ChatRoom")
	}
	if room.PetSitterID != petSitterID {
		return exceptions.NewUnauthorizedError("You are not a participant of this chat room")
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
		if room.UserLastReadMessageID != nil && request.LastReadMessageID <= *room.UserLastReadMessageID {
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
