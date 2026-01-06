package postgres

import (
	"hona/backend/internal/domain/entities"
	"log"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) CreateRoom(room *entities.ChatRoom) error {
	return cr.db.Create(&room).Error
}

func (cr *ChatRepository) GetUserAndPetSitterRoom(userID uint, petSitterID uint) (*entities.ChatRoom, error) {
	var foundRoom entities.ChatRoom
	err := cr.db.Where("(user_id = ? AND pet_sitter_id = ?)", userID, petSitterID).First(&foundRoom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println("Error retrieving room:", err)
		return nil, err
	}
	log.Println("Found room:", foundRoom)
	return &foundRoom, nil
}

func (cr *ChatRepository) GetAllRoomsByUserID(userID uint) ([]*entities.ChatRoom, error) {
	var rooms []*entities.ChatRoom
	err := cr.db.Where("user_id = ?", userID).Find(&rooms).Error
	return rooms, err
}

func (cr *ChatRepository) GetAllRoomsByPetSitterID(petSitterID uint) ([]*entities.ChatRoom, error) {
	var rooms []*entities.ChatRoom
	err := cr.db.Where("pet_sitter_id = ?", petSitterID).Find(&rooms).Error
	return rooms, err
}

func (cr *ChatRepository) UpdateRoomStatus(roomID uint, status uint) error {
	return cr.db.Model(&entities.ChatRoom{}).Where("id = ?", roomID).Update("status", status).Error
}

func (cr *ChatRepository) UpdateRoomBlockedBy(roomID uint, blockedBy string) error {
	return cr.db.Model(&entities.ChatRoom{}).Where("id = ?", roomID).Update("blocked_by", blockedBy).Error
}

func (cr *ChatRepository) GetRoomByID(roomID uint) (*entities.ChatRoom, error) {
	var room entities.ChatRoom
	err := cr.db.First(&room, roomID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (cr *ChatRepository) GetRequestIDByRoomID(roomID uint) (*uint, error) {
	var room entities.ChatRoom
	err := cr.db.Select("request_id").First(&room, roomID).Error
	if err != nil {
		return nil, err
	}
	return room.RequestID, nil
}

func (cr *ChatRepository) CreateMessage(message *entities.ChatMessage) error {
	return cr.db.Create(message).Error
}

func (cr *ChatRepository) UpdateRoom(room *entities.ChatRoom) error {
	return cr.db.Save(room).Error
}

func (cr *ChatRepository) GetAllRooms(senderID uint, options *QueryOptions) ([]*entities.ChatRoom, int64, error) {
	var rooms []*entities.ChatRoom
	var totalCount int64
	log.Println("Getting all rooms for sender ID:", senderID)
	query := cr.db.Model(&entities.ChatRoom{}).Where("user_id = ? OR pet_sitter_id = ?", senderID, senderID)
	log.Println("Base query constructed")
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	log.Println("Total rooms found:", totalCount)
	query = query.Order("last_message_id DESC NULLS LAST")
	if options.Pagination != nil {
		paginationModifier := NewPaginationModifier(options.Pagination.Offset, options.Pagination.Limit)
		query = paginationModifier.Apply(query)
	}
	if err := query.Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, totalCount, nil
}

func (cr *ChatRepository) UnreadMessageCount(roomID uint, senderID uint, lastReadMessageID *uint) (int64, error) {
	var count int64
	query := cr.db.Model(&entities.ChatMessage{}).
		Where("room_id = ?", roomID).
		Where("sender_id != ?", senderID).
		Where("is_deleted = ?", false)
	if lastReadMessageID != nil {
		query = query.Where("id > ?", *lastReadMessageID)
	}
	err := query.Count(&count).Error
	return count, err
}

func (cr *ChatRepository) FindLastMessageByID(messageID *uint) (*entities.ChatMessage, error) {
	if messageID == nil {
		return nil, nil
	}
	var message entities.ChatMessage
	err := cr.db.First(&message, messageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (cr *ChatRepository) FindLastActiveMessageInRoom(roomID uint) (*entities.ChatMessage, error) {
	var message entities.ChatMessage
	err := cr.db.
		Where("room_id = ? AND is_deleted = ?", roomID, false).
		Order("id DESC").
		First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (cr *ChatRepository) GetMessagesByRoomID(roomID uint, options *QueryOptions) ([]*entities.ChatMessage, int64, error) {
	var messages []*entities.ChatMessage
	var totalCount int64
	query := cr.db.Model(&entities.ChatMessage{}).Where("room_id = ?", roomID)
	if options.Filters != nil {
		filterModifier := NewFilterModifier(options.Filters.Filters)
		query = filterModifier.Apply(query)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if options.Sorting != nil && len(options.Sorting.Sorts) > 0 {
		sortModifier := NewSortModifier(options.Sorting.Sorts)
		query = sortModifier.Apply(query)
	} else {
		query = query.Order("id DESC")
	}
	if options.Pagination != nil {
		paginationModifier := NewPaginationModifier(options.Pagination.Offset, options.Pagination.Limit)
		query = paginationModifier.Apply(query)
	}
	if err := query.Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, totalCount, nil
}

func (cr *ChatRepository) FindMessageByID(messageID uint) (*entities.ChatMessage, error) {
	var message entities.ChatMessage
	err := cr.db.First(&message, messageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}

func (cr *ChatRepository) FindReaction(messageID uint, userID uint, emoji string) (*entities.ChatMessageReaction, error) {
	var reaction entities.ChatMessageReaction
	err := cr.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).First(&reaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func (cr *ChatRepository) AddReaction(reaction *entities.ChatMessageReaction) error {
	return cr.db.Create(reaction).Error
}

func (cr *ChatRepository) RemoveReaction(messageID uint, userID uint, emoji string) error {
	return cr.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Delete(&entities.ChatMessageReaction{}).Error
}

func (cr *ChatRepository) GetReactionsByMessageIDs(messageIDs []uint) (map[uint][]*entities.ChatMessageReaction, error) {
	result := make(map[uint][]*entities.ChatMessageReaction)
	if len(messageIDs) == 0 {
		return result, nil
	}
	var reactions []*entities.ChatMessageReaction
	if err := cr.db.Where("message_id IN ?", messageIDs).Find(&reactions).Error; err != nil {
		return nil, err
	}
	for _, reaction := range reactions {
		result[reaction.MessageID] = append(result[reaction.MessageID], reaction)
	}
	return result, nil
}

func (cr *ChatRepository) UpdateMessage(message *entities.ChatMessage) error {
	return cr.db.Save(message).Error
}
func (cr *ChatRepository) DeleteMessageByID(messageID uint) error {
	return cr.db.Delete(&entities.ChatMessage{}, messageID).Error
}
