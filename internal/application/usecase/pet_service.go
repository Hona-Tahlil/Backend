package usecase

import (
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/domain/entities"
)

type PetService interface {
	AddPet(info pet.AddPetRequest) error
	UpdatePet(info pet.UpdatePetRequest) error
	RemovePet(info pet.RemovePetRequest) error
	GetPetsBasicData(info pet.GetPetsBasicDataRequest) ([]pet.PetBasicDataResponse, error)
	GetPetFullData(info pet.GetPetFullDataRequest) (*pet.PetFullDataResponse, error)
	GetAllPetKinds() []pet.PetKindResponse
	GetPetKindSpecies(info pet.GetPetKindSpecies) []pet.PetSpeciesResponse
	FindPetByID(id uint) (*entities.Pet, error)
	GetPetsBasicDataResponse(pets []entities.Pet) ([]pet.PetBasicDataResponse, error)
	GetPetsInUser(userPets []entities.Pet, petIDs []uint) ([]entities.Pet, error)
	GetPetNames(pets []entities.Pet) []string
	FindUserPetsByID(userID uint) ([]entities.Pet, error)
}
