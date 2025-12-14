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

func (pr *PetSitterRepository) FindPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	var foundPetsitter entities.PetSitter

	if result := pr.db.First(&foundPetsitter, "userID = ?", id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundPetsitter, nil
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
func (pr *PetSitterRepository) PreloadServices(petSitter *entities.PetSitter) error {
	return pr.db.Preload("Services").First(petSitter, petSitter.ID).Error
}

func (pr *PetSitterRepository) CreatePetSitter(petSitter *entities.PetSitter) error {
	return pr.db.Create(petSitter).Error
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

func (pr *PetSitterRepository) FindServiceByID(id uint) (*entities.Service, error) {
	var foundService entities.Service

	if result := pr.db.First(&foundService, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &foundService, nil
}
