package service

import (
	"errors"
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainstorage "hona/backend/internal/domain/storage"
)

type PetSitterService struct {
	unitOfWork  ports.UnitOfWork
	storage     domainstorage.Storage
	userService usecase.UserService
}

func NewPetSitterService(unitOfWork ports.UnitOfWork, storage domainstorage.Storage, userService usecase.UserService) *PetSitterService {
	return &PetSitterService{
		unitOfWork:  unitOfWork,
		storage:     storage,
		userService: userService,
	}
}

func (ps *PetSitterService) CreateSignupSession(PetsitterInfo petsitter.GetPetSitterRequest) (*petsitter.PetSitterStatusResponse, error) {
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

func (ps *PetSitterService) SubmitPersonalInfo(petSitterInfo petsitter.SubmitPersonalInfoRequest) (*petsitter.PetSitterStatusResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundUser, err := userRepo.FindUserByID(petSitterInfo.UserID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(petSitterInfo.UserID)
	if err != nil {
		return nil, err
	}
	if foundPetSitter.OnboardingStep == enums.OBS_Done {
		return &petsitter.PetSitterStatusResponse{
			UserID:         foundPetSitter.UserID,
			Status:         foundPetSitter.Status,
			OnboardingStep: foundPetSitter.OnboardingStep,
		}, err
	}
	if foundPetSitter.Status == enums.PSS_Rejected || foundPetSitter.Status == enums.PSS_Suspended {
		return &petsitter.PetSitterStatusResponse{
			UserID:         foundPetSitter.UserID,
			Status:         foundPetSitter.Status,
			OnboardingStep: foundPetSitter.OnboardingStep,
		}, err
	}
	foundUser.FirstName = petSitterInfo.FirstName
	foundUser.LastName = petSitterInfo.LastName
	foundUser.Email = petSitterInfo.Email
	foundUser.Phone = &petSitterInfo.PhoneNumber
	foundUser.FirstName = petSitterInfo.FirstName
	foundUser.FirstName = petSitterInfo.FirstName

	err = userRepo.UpdateUser(foundUser)
	err = petSitterRepo.UpdatePetSitter(foundPetSitter)
	if err != nil {
		return nil, err
	}

	return &petsitter.PetSitterStatusResponse{
		UserID:         foundPetSitter.UserID,
		Status:         foundPetSitter.Status,
		OnboardingStep: foundPetSitter.OnboardingStep,
	}, nil

}

func (ps *PetSitterService) GetPersonalInfo(userID uint) (*petsitter.PersonalInfoResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	addressRepo := ps.unitOfWork.Factory().AddressRepository()
	foundUser, err := userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(userID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadAddress(foundUser)
	if err != nil {
		return nil, err
	}
	err = addressRepo.PreloadProvince(&foundUser.Address)
	if err != nil {
		return nil, err
	}
	err = addressRepo.PreloadCity(&foundUser.Address)
	if err != nil {
		return nil, err
	}
	return &petsitter.PersonalInfoResponse{
		UserID:         foundUser.ID,
		FirstName:      foundUser.FirstName,
		LastName:       foundUser.LastName,
		Email:          foundUser.Email,
		PhoneNumber:    *foundUser.Phone,
		Gender:         foundUser.Gender,
		BirthDate:      foundUser.BirthDate,
		Province:       foundUser.Address.Province.Name,
		City:           foundUser.Address.City.Name,
		Address:        foundUser.Address.StreetAddress,
		Pelak:          foundUser.Address.HouseNumber,
		Vahed:          foundUser.Address.Unit,
		PostalCode:     *foundUser.Address.PostalCode,
		Status:         foundPetSitter.Status,
		OnboardingStep: foundPetSitter.OnboardingStep,
	}, nil

}

func (ps *PetSitterService) UploadDocuments(info petsitter.UploadDocumentsRequest) (*petsitter.PetSitterStatusResponse, error) {
	var profileKey *string
	if info.File != nil {
		profileKeyValue := ps.getStorageKey(info.UserID)
		profileKey = &profileKeyValue
		if err := ps.storage.UploadFile(enums.PetSitterFile, *profileKey, info.File); err != nil {
			return nil, err
		}
	}

}

func (ps *PetSitterService) GetDocuments(userID uint) (*petsitter.DocumentResponse, error) {
	foundPetSitter, err := ps.findPetSitterByID(userID)
	if err != nil {
		return nil, err
	}
	var profileKey *string
	profileKeyValue := ps.getStorageKey(info.UserID)
	profileKey = &profileKeyValue
	//GetPresignedURL(bucketType enums.BucketType, key string, expiration time.Duration) (string, error)
	fileurl ,_ := ps.storage.GetPresignedURL(enums.PetSitterFile, *profileKey, time.Duration)
	return &petsitter.DocumentResponse{
		UserID: userID,
		File: fileurl,
	}, nil
}

func (ps *PetSitterService) SubmitSkills(SkillsInfo petsitter.SubmitSkillsRequest) (*petsitter.PetSitterStatusResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundUser, err := userRepo.FindUserByID(SkillsInfo.UserID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	if foundUser.PetSitter == nil {
		return nil, errors.New("petsitter record missing")
	}
	foundPetSitter, err := ps.findPetSitterByID(SkillsInfo.UserID)
	if err != nil {
		return nil, err
	}
	err = petSitterRepo.PreloadServices(foundPetSitter)
	if err != nil {
		return nil, err
	}
	services := ps.GetServicesResponse(SkillsInfo.Services)
	if foundPetSitter.OnboardingStep != enums.OBS_Documents {
		return nil, errors.New("invalid onboarding step: cannot submit skills now")
	}
	foundPetSitter.Bio = &SkillsInfo.Bio
	foundPetSitter.Services = services
	for _, petkind := range foundPetSitter.PetKinds {
		foundPetSitter.PetKinds = append(foundPetSitter.PetKinds, petkind)
	}
	foundPetSitter.OnboardingStep = enums.OBS_Done

	err = petSitterRepo.UpdatePetSitter(foundPetSitter)
	if err != nil {
		return nil, err
	}
	return &petsitter.PetSitterStatusResponse{
		UserID:         foundPetSitter.UserID,
		Status:         foundPetSitter.Status,
		OnboardingStep: foundPetSitter.OnboardingStep,
	}, nil

}

func (ps *PetSitterService) GetPetsitterStatus(userID uint) (*petsitter.PetSitterStatusResponse, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(userID)
	if err != nil {
		return nil, err
	}
	return &petsitter.PetSitterStatusResponse{
		UserID:         foundPetSitter.UserID,
		OnboardingStep: foundPetSitter.OnboardingStep,
		Status:         foundPetSitter.Status,
	}, nil
}

func (ps *PetSitterService) getStorageKey(userID uint) string {
	return "file-" + string(userID) + "-" + fmt.Sprint(userID)
}

func (ps *PetSitterService) findPetSitterByID(id uint) (*entities.PetSitter, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(id)
	if err != nil {
		return nil, err
	}
	if foundPetSitter == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
		return nil, NotFoundError
	}
	err = petSitterRepo.PreloadServices(foundPetSitter)
	if err != nil {
		return nil, err
	}
	return foundPetSitter, nil
}


func (ps *PetSitterService) GetServicesResponse(Services []enums.ServiceType) []entities.Service {
	r := make([]entities.Service, 0)
	for _, service := range Services {
		r = append(r, entities.Service{
			Type: service,
		})
	}
	return r
}

