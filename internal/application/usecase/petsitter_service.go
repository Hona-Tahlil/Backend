package usecase

type PetSitterService interface {
	GetAllPetSitters(page, count int) (interface{}, error)
}