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
	return pr.db.Preload("schedule").First(petSitter, petSitter.ID).Error
}

func (pr *PetSitterRepository) CreatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Create(petSitter).Error
}

func (pr *PetSitterRepository) SavePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Save(petSitter).Error
}