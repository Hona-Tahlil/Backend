package domainpostgres

import "hona/backend/internal/domain/entities"

type PetSitterRepository interface {
<<<<<<< HEAD
	CreatePetSitter(petSitter *entities.PetSitter) error
	UpdatePetSitter(petSitter *entities.PetSitter) error
	PreloadServices(petSitter *entities.PetSitter) error
	FindPetSitterByUserID(id uint) (*entities.PetSitter, error)
	GetAllPetSitters(limit, offset int) ([]entities.PetSitter, error)
	GetPetSittersCount() (int64, error)
=======
	PreloadFields(petSitter *entities.PetSitter, fields []string) error
	EditPetSitter(petSitter *entities.PetSitter) error
	FindPetSitterByID(id uint) (*entities.PetSitter, error)
>>>>>>> dev
}
