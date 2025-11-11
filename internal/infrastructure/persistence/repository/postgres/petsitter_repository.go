package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type PetSitterRepository struct {
	db *gorm.DB
}

func NewPetSitterRepository(db *gorm.DB) *PetSitterRepository {
	return &PetSitterRepository{
		db: db,
	}
}

func (pr *PetSitterRepository) PreloadSchedule(petSitter *entities.PetSitter) error {
	return pr.db.Preload("Schedule").First(petSitter, petSitter.ID).Error
}
