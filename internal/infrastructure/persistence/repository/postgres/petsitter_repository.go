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

func (pr *PetSitterRepository) PreloadServices(petSitter *entities.PetSitter) error {
	return pr.db.Preload("Services").First(petSitter, petSitter.ID).Error
}

func (pr *PetSitterRepository) CreatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Create(petSitter).Error
}

func (pr *PetSitterRepository) UpdatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Save(petSitter).Error
}

func (pr *PetSitterRepository) FindPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	var petsitter entities.PetSitter
	err := pr.db.First(&petsitter, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// return &	, nil

	// if result := pr.db.First(&petsitter, "userid = ?", id); result.Error != nil {
	// 	if result.Error == gorm.ErrRecordNotFound {
	// 		return nil, nil
	// 	}
	// 	return nil, result.Error
	// }
	return &petsitter, nil
}

func (pr *PetSitterRepository) GetAllPetSitters(limit, offset int) ([]entities.PetSitter, error) {
	var petSitters []entities.PetSitter
	err := pr.db.Limit(limit).Offset(offset).Find(&petSitters).Error
	if err != nil {
		return nil, err
	}
	return petSitters, nil
}

func (pr *PetSitterRepository) GetPetSittersCount() (int64, error) {
	var count int64
	err := pr.db.Model(&entities.PetSitter{}).Count(&count).Error
	return count, err
}
