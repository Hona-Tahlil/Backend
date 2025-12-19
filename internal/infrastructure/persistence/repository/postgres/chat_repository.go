package postgres

import (
	"hona/backend/internal/domain/entities"

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

func (cr *ChatRepository) SaveMessage(message *entities.ChatMessage) error {
	return cr.db.Create(message).Error
}
