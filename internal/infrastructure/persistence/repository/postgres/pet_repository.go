package postgres

import (
	"hona/backend/internal/domain/entities"

	"gorm.io/gorm"
)

type PetRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{
		db: db,
	}
}

func (pr *PetRepository) CreatePet(pet *entities.Pet) error {
	return pr.db.Create(pet).Error
}
