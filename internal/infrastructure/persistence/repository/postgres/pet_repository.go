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

func (pr *PetRepository) FindPet(name string, userID uint) (*entities.Pet, error) {
	var pet entities.Pet
	err := pr.db.Where("name = ? AND user_id = ?", name, userID).First(&pet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

func (pr *PetRepository) FindPetByID(id uint) (*entities.Pet, error) {
	var pet entities.Pet
	err := pr.db.First(&pet, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

func (pr *PetRepository) UpdatePet(pet *entities.Pet) error {
	return pr.db.Save(pet).Error
}

func (pr *PetRepository) RemovePet(pet *entities.Pet) error {
	return pr.db.Delete(pet).Error
}

func (pr *PetRepository) PreloadUserPets(user *entities.User) error {
	return pr.db.Preload("Pets").First(user, user.ID).Error
}

func (pr *PetRepository) FindUserPetsByID(userID uint) ([]entities.Pet, error) {
	var pets []entities.Pet
	err := pr.db.Where("userID = ? AND type = ?", userID, "regular").Find(&pets).Error
	if err != nil {
		return nil, err
	}
	return pets, nil
}

func (pr *PetRepository) FindRequestPetsByID(requestID uint) ([]entities.Pet, error) {
	var pets []entities.Pet
	err := pr.db.Where("requestID = ? AND type = ?", requestID, "request").Find(&pets).Error
	if err != nil {
		return nil, err
	}
	return pets, nil
}
