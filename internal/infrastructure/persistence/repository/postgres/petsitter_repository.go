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

func (pr *PetSitterRepository) FindPetSitterByID(id uint) (*entities.PetSitter, error) {
	var petSitter entities.PetSitter
	result := pr.db.First(&petSitter, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &petSitter, nil
}
