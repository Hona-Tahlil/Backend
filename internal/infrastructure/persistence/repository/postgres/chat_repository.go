package postgres

import (
	"hona/backend/internal/domain/entities"
	domainpostgres "hona/backend/internal/domain/ports/postgres"

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
		return nil, err
	}
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

func (cr *ChatRepository) GetRequestIDByRoomID(roomID uint) (uint, error) {
	var room entities.ChatRoom
	err := cr.db.Select("request_id").First(&room, roomID).Error
	if err != nil {
		return 0, err
	}
	return room.RequestID, nil
}

func (cr *ChatRepository) CreateMessage(message *entities.ChatMessage) error {
	return cr.db.Create(message).Error
}

func (cr *ChatRepository) UpdateRoom(room *entities.ChatRoom) error {
	return cr.db.Save(room).Error
}

func (cr *ChatRepository) GetAllRooms(senderID uint, options *domainpostgres.QueryOptions) ([]*entities.ChatRoom, int64, error) {
	var rooms []*entities.ChatRoom
	var totalCount int64
	query := cr.db.Model(&entities.ChatRoom{}).Where("user_id = ? OR pet_sitter_id = ?", senderID, senderID)
	if options.Filters != nil {
		filterModifier := NewFilterModifier(options.Filters.Filters)
		query = filterModifier.Apply(query)
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	if options.Sorting != nil {
		sortModifier := NewSortModifier(options.Sorting.Sorts)
		query = sortModifier.Apply(query)
	}
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
		Where("room_id = ? AND id > ?", roomID, lastReadMessageID).
		Where("sender_id != ?", senderID).
		Count(&count).Error
	return count, query
}

func (cr *ChatRepository) FindLastMessageByID(messageID *uint) (*entities.ChatMessage, error) {
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

func (cr *ChatRepository) GetMessagesByRoomID(roomID uint, options *domainpostgres.QueryOptions) ([]*entities.ChatMessage, int64, error) {
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
	if options.Sorting != nil {
		sortModifier := NewSortModifier(options.Sorting.Sorts)
		query = sortModifier.Apply(query)
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

func (cr *ChatRepository) UpdateMessage(message *entities.ChatMessage) error {
	return cr.db.Save(message).Error
}
func (cr *ChatRepository) DeleteMessageByID(messageID uint) error {
	return cr.db.Delete(&entities.ChatMessage{}, messageID).Error
}