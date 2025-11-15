package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
)

type PetSitterService struct {
	unitOfWork ports.UnitOfWork
}

func NewPetSitterService(unitOfWork ports.UnitOfWork) *PetSitterService {
	return &PetSitterService{
		unitOfWork: unitOfWork,
	}
}
func (ps *PetSitterService) CreateSignupSession(PetsitterInfo *petsitter.GetPetSitterRequest) (*petsitter.PetSitterStatusResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundUser, err := userRepo.FindUserByID(PetsitterInfo.UserID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	if foundUser.PetSitter != nil {
		return &petsitter.PetSitterStatusResponse{
			UserID:         foundUser.ID,
			Status:         foundUser.PetSitter.Status,
			OnboardingStep: foundUser.PetSitter.OnboardingStep,
		}, nil
	}
	newPetSitter := &entities.PetSitter{
		UserID:         foundUser.ID,
		Status:         enums.PSS_Draft,
		OnboardingStep: enums.OBS_Review,
	}
	err = petSitterRepo.CreatePetSitter(newPetSitter)
	if err != nil {
		return nil, err
	}
	return &petsitter.PetSitterStatusResponse{
		UserID:         foundUser.ID,
		Status:         foundUser.PetSitter.Status,
		OnboardingStep: foundUser.PetSitter.OnboardingStep,
	}, nil
}

func (ps *PetSitterService) GetAndUpdateUser(PetsitterInfo petsitter.GetPetSitterRequest) (*petsitter.PetSitterGetResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()

	foundUser, err := userRepo.FindUserByID(PetsitterInfo.UserID)
	if foundUser == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
		return nil, NotFoundError
	}
	if err != nil {
		return nil, err
	}

	return &petsitter.PetSitterGetResponse{}, nil

}

func (ps *PetSitterService) SubmitPersonalInfo(petsitterInfo petsitter.FirstSubmit) error {
	
}