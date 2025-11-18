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
	foundUser, err := ps.userService.FindVerifiedUserByID(PetsitterInfo.UserID)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	if foundUser.PetSitter != nil {
		return &petsitter.PetSitterStatusResponse{
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
		Status:         newPetSitter.Status,
		OnboardingStep: newPetSitter.OnboardingStep,
	}, nil
}

func (ps *PetSitterService) SubmitPersonalInfo(petSitterInfo petsitter.SubmitPersonalInfoRequest) error {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	addressRepo := ps.unitOfWork.Factory().AddressRepository()
	foundUser, err := ps.userService.FindVerifiedUserByID(petSitterInfo.UserID)
	if err != nil {
		return err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return err
	}
	if foundUser.PetSitter == nil {
		return err
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Profile)
	if err != nil {
		return err
	}
	err = userRepo.PreloadAddress(foundUser)
	if err != nil {
		return err
	}
	err = addressRepo.PreloadProvince(&foundUser.Address)
	if err != nil {
		return err
	}
	err = addressRepo.PreloadCity(&foundUser.Address)
	if err != nil {
		return err
	}
	foundUser.FirstName = petSitterInfo.FirstName
	foundUser.LastName = petSitterInfo.LastName
	foundUser.Email = petSitterInfo.Email
	foundUser.Gender = petSitterInfo.Gender
	foundUser.BirthDate = petSitterInfo.BirthDate
	foundUser.Phone = &petSitterInfo.Phone
	foundUser.Address.StreetAddress = petSitterInfo.Address
	foundUser.Address.Province.Name = petSitterInfo.Province
	foundUser.Address.City.Name = petSitterInfo.City
	foundUser.Address.HouseNumber = petSitterInfo.HouseNumber
	foundUser.Address.Unit = petSitterInfo.Unit
	foundUser.PetSitter.Status = enums.PSS_Draft
	foundUser.PetSitter.OnboardingStep = enums.OBS_Profile
	err = userRepo.UpdateUser(foundUser)
	if err != nil {
		return err
	}
	return nil

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
		FirstName:      foundUser.FirstName,
		LastName:       foundUser.LastName,
		Email:          foundUser.Email,
		PhoneNumber:    *foundUser.Phone,
		Gender:         foundUser.Gender,
		BirthDate:      foundUser.BirthDate,
		Province:       foundUser.Address.Province.Name,
		City:           foundUser.Address.City.Name,
		Address:        foundUser.Address.StreetAddress,
		HouseNumber:    foundUser.Address.HouseNumber,
		Unit:           foundUser.Address.Unit,
		PostalCode:     *foundUser.Address.PostalCode,
		Status:         foundPetSitter.Status,
		OnboardingStep: foundPetSitter.OnboardingStep,
	}, nil

}

func (ps *PetSitterService) UploadDocuments(info petsitter.UploadDocumentsRequest) error {
	foundPetSitter, err := ps.FindPetSitterByID(info.UserID)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStatus(foundPetSitter.Status)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStep(foundPetSitter.OnboardingStep, enums.OBS_Documents)
	if err != nil {
		return err
	}
	var fileKey *string
	CertificateKeys := make([]string, 0)
	if info.CertificateFiles != nil {
		for _, file := range info.CertificateFiles {
			fileKeyValue := ps.getStorageKey(info.UserID)
			fileKey = &fileKeyValue
			if err := ps.storage.UploadFile(enums.PetSitterFile, *fileKey, file); err != nil {
				return err
			}
			CertificateKeys = append(CertificateKeys, *fileKey)
		}
	}
	FileKeys := make([]string, 0)
	if info.Files != nil {
		for _, file := range info.Files {
			fileKeyValue := ps.getStorageKey(info.UserID)
			fileKey = &fileKeyValue
			if err := ps.storage.UploadFile(enums.PetSitterFile, *fileKey, file); err != nil {
				return err
			}
			FileKeys = append(FileKeys, *fileKey)
		}
	}
	foundPetSitter.CertificateKeys = CertificateKeys
	foundPetSitter.FileKeys = FileKeys
	foundPetSitter.OnboardingStep = enums.OBS_Documents
	foundPetSitter.Status = enums.PSS_InReview
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.UpdatePetSitter(foundPetSitter)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PetSitterService) GetDocuments(userID uint) (*petsitter.DocumentResponse, error) {
	foundPetSitter, err := ps.FindPetSitterByID(userID)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStatus(foundPetSitter.Status)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStep(foundPetSitter.OnboardingStep, enums.OBS_Documents)
	if err != nil {
		return nil, err
	}
	CertificateFiles := make([]string, 0)
	for _, certKey := range foundPetSitter.CertificateKeys {
		certurl, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, certKey, 2)
		CertificateFiles = append(CertificateFiles, certurl)
	}
	Files := make([]string, 0)
	for _, fileKey := range foundPetSitter.FileKeys {
		fileurl, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, fileKey, 2)
		Files = append(Files, fileurl)
	}
	return &petsitter.DocumentResponse{
		CertificateFiles: CertificateFiles,
		Files:            Files,
	}, nil
}

func (ps *PetSitterService) SubmitSkills(SkillsInfo petsitter.SubmitSkillsRequest) error {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundUser, err := userRepo.FindUserByID(SkillsInfo.UserID)
	if err != nil {
		return err
	}
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return err
	}
	if foundUser.PetSitter == nil {
		return errors.New("petsitter record missing")
	}
	foundPetSitter, err := ps.FindPetSitterByID(SkillsInfo.UserID)
	if err != nil {
		return err
	}
	err = petSitterRepo.PreloadServices(foundPetSitter)
	if err != nil {
		return err
	}
	services := ps.GetServicesResponse(SkillsInfo.Services)
	if foundPetSitter.OnboardingStep != enums.OBS_Documents {
		return errors.New("invalid onboarding step: cannot submit skills now")
	}
	foundPetSitter.Bio = &SkillsInfo.Bio
	foundPetSitter.Services = services
	for _, petkind := range foundPetSitter.PetKinds {
		foundPetSitter.PetKinds = append(foundPetSitter.PetKinds, petkind)
	}
	foundPetSitter.OnboardingStep = enums.OBS_Done

	err = petSitterRepo.UpdatePetSitter(foundPetSitter)
	if err != nil {
		return err
	}
	return nil

}

func (ps *PetSitterService) GetPetsitterStatus(userID uint) (*petsitter.PetSitterStatusResponse, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(userID)
	if err != nil {
		return nil, err
	}
	return &petsitter.PetSitterStatusResponse{
		OnboardingStep: foundPetSitter.OnboardingStep,
		Status:         foundPetSitter.Status,
	}, nil
}

func (ps *PetSitterService) getStorageKey(userID uint) string {
	return "file-" + fmt.Sprint(userID)
}

func (ps *PetSitterService) FindPetSitterByID(id uint) (*entities.PetSitter, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(id)
	if err != nil {
		return nil, err
	}
	if foundPetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
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

func (ps *PetSitterService) CheckPetSitterStatus(pss enums.PetSitterStatus) error {
	if pss == enums.PSS_Rejected {
		return errors.New("operation not allowed: petsitter is rejected")
	}
	if pss == enums.PSS_Suspended {
		return errors.New("operation not allowed: petsitter is suspended")
	}
	if pss == enums.PSS_Active {
		return errors.New("operation not allowed: petsitter is active")
	}
	if pss == enums.PSS_InReview {
		return errors.New("operation not allowed: petsitter is in review")
	}
	return nil
}

func (ps *PetSitterService) CheckPetSitterStep(currentStep enums.OnboardingStep, requiredStep enums.OnboardingStep) error {
	if (currentStep + 1) < requiredStep {
		return errors.New("operation not allowed: invalid onboarding step")
	}
	return nil
}
