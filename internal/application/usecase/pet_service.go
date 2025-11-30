package usecase

import "hona/backend/internal/application/dto/pet"

type PetService interface {
	AddPet(info pet.AddPetRequest) error
	UpdatePet(info pet.UpdatePetRequest) error
	RemovePet(info pet.RemovePetRequest) error
	GetPetsBasicData(info pet.GetPetsBasicDataRequest) ([]pet.PetBasicDataResponse, error)
	GetPetFullData(info pet.GetPetFullDataRequest) (*pet.PetFullDataResponse, error)
	GetAllPetKinds() []pet.PetKindResponse
	GetPetKindSpecies(info pet.GetPetKindSpecies) []pet.PetSpeciesResponse
}
