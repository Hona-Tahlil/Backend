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

func (pr *PetSitterRepository) PreloadFields(petSitter *entities.PetSitter, fields []string) error {
	for _, field := range fields {
		err := pr.db.Preload(field).First(petSitter, petSitter.ID).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (pr *PetSitterRepository) EditPetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Save(petSitter).Error
}
