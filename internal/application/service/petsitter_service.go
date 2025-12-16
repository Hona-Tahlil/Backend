package service

import (
	"errors"
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainstorage "hona/backend/internal/domain/storage"

	"github.com/samber/lo"
)

type PetSitterService struct {
	unitOfWork     ports.UnitOfWork
	storage        domainstorage.Storage
	userService    usecase.UserService
	addressService usecase.AddressService
}

func NewPetSitterService(unitOfWork ports.UnitOfWork, storage domainstorage.Storage, userService usecase.UserService, addressService usecase.AddressService) *PetSitterService {
	return &PetSitterService{
		unitOfWork:     unitOfWork,
		storage:        storage,
		userService:    userService,
		addressService: addressService,
	}
}

func (ps *PetSitterService) GetPetSitterFreeSlotsResponse(petSitter *entities.PetSitter) ([]calendarslot.CalendarSlotInfoResponse, error) {
	if petSitter.Schedule == nil {
		return nil, nil
	}
	petSitterCalendarSlots := petSitter.Schedule
	freeSlots := make([]entities.CalendarSlot, len(petSitterCalendarSlots))
	for i, slot := range petSitterCalendarSlots {
		if slot.Status != enums.Free {
			continue
		}
		freeSlots[i] = slot
	}

	return ps.GetCalendarSlotsResponse(freeSlots), nil
}

func (ps *PetSitterService) GetServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, len(petSitter.Services))
	if petSitter.Services == nil {
		return r, nil
	}
	for i, service := range petSitter.Services {
		r[i] = ps.GetServiceResponse(&service)
	}
	return r, nil
}

func (ps *PetSitterService) GetAvailableServicesResponse(petSitter *entities.PetSitter) ([]servicedto.ServiceInfoResponse, error) {
	r := make([]servicedto.ServiceInfoResponse, len(petSitter.Services))
	if petSitter.Services == nil {
		return r, nil
	}
	for i, service := range petSitter.Services {
		if service.Price != 0 {
			r[i] = ps.GetServiceResponse(&service)
		}
	}
	return r, nil
}

func (ps *PetSitterService) GetPetSitterByUserID(id uint) (*entities.PetSitter, error) {
	user, err := ps.userService.FindUserByID(id)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = ps.userService.PreloadFields(user, []string{"PetSitter"})
	if err != nil {
		return nil, err
	}
	if user.PetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if user.PetSitter.Status != enums.PSS_Active {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return user.PetSitter, nil
}

func (ps *PetSitterService) GetPetSitterByID(id uint) (*entities.PetSitter, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByID(id)
	if err != nil {
		return nil, err
	}
	if foundPetSitter == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return foundPetSitter, nil
}

func (ps *PetSitterService) PreloadFields(petSitter *entities.PetSitter, fields []string) error {
	return ps.unitOfWork.Factory().PetSitterRepository().PreloadFields(petSitter, fields)
}

func (ps *PetSitterService) ValidatePets(pets []entities.Pet, petKinds []enums.PetKind) error {
	for _, pet := range pets {
		flag := false
		for _, kind := range petKinds {
			if kind == pet.Kind {
				flag = true
				break
			}
		}
		if !flag {
			var ve exceptions.ValidationErrors
			ve.AddError(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			return &ve
		}
	}

	return nil
}

func (ps *PetSitterService) ValidateService(services []entities.Service, serviceID uint) (*entities.Service, error) {
	for _, service := range services {
		if service.ID == serviceID && service.Price != 0 {
			return &entities.Service{
				PetSitterID: service.PetSitterID,
				Type:        service.Type,
				Price:       service.Price,
				Description: service.Description,
				Kind:        bootstrap.Run().Constants.EntityConstants.Request,
			}, nil
		}
	}
	var ve exceptions.ValidationErrors
	ve.AddError(bootstrap.Run().Constants.ErrorFields.Service, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
	return nil, &ve
}

// AutoUpdateSlots can be accept or cancel
func (ps *PetSitterService) AutoUpdateSlots(petSitter *entities.PetSitter, calendarSlots []entities.CalendarSlot, accept bool) error {
	newSlots := make([]entities.CalendarSlot, 0)
	petSitterSlots := petSitter.Schedule
	for _, petSitterSlot := range petSitterSlots {
		for _, calendarSlot := range calendarSlots {
			if petSitterSlot.Date.Equal(calendarSlot.Date) {
				if accept && petSitterSlot.Status == enums.Free {
					petSitterSlot.Slots, _ = lo.Difference(petSitterSlot.Slots, calendarSlot.Slots)
				} else if !accept && petSitterSlot.Status == enums.Booked {
					petSitterSlot.Slots, _ = lo.Difference(petSitterSlot.Slots, calendarSlot.Slots)
				}
				newSlots = append(newSlots, ps.makeSlot(&calendarSlot, accept))
			}
		}
	}
	petSitter.Schedule = append(petSitter.Schedule, newSlots...)
	ps.removeEmptySitterSlots(petSitter)
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.EditPetSitter(petSitter)
}

func (ps *PetSitterService) removeEmptySitterSlots(petSitter *entities.PetSitter) {
	for i, slot := range petSitter.Schedule {
		if len(slot.Slots) == 0 {
			petSitter.Schedule = append(petSitter.Schedule[:i], petSitter.Schedule[i+1:]...)
		}
	}
}

func (ps *PetSitterService) makeSlot(calendarSlot *entities.CalendarSlot, accept bool) entities.CalendarSlot {
	var slot entities.CalendarSlot
	slot.Date = calendarSlot.Date
	if accept {
		slot.Status = enums.Booked
	} else {
		slot.Status = enums.Free
	}
	slot.Slots = calendarSlot.Slots
	return slot
}

func (ps *PetSitterService) CreateSignupSession(PetsitterInfo petsitter.GetPetSitterRequest) (*petsitter.PetSitterStatusResponse, error) {
	userRepo := ps.unitOfWork.Factory().UserRepository()
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundUser, err := ps.userService.FindVerifiedUserByID(PetsitterInfo.UserID)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
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

func (ps *PetSitterService) GetPersonalInfo(userID uint) (*petsitter.PersonalInfoResponse, error) {
	foundUser, err := ps.userService.FindVerifiedUserByID(userID)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, err
	}
	err = userRepo.PreloadAddress(foundUser)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Documents)
	if err != nil {
		return nil, err
	}
	r := &petsitter.PersonalInfoResponse{
		FirstName:      foundUser.FirstName,
		LastName:       foundUser.LastName,
		Email:          foundUser.Email,
		PhoneNumber:    *foundUser.Phone,
		Gender:         foundUser.Gender,
		Province:       foundUser.Address.Province,
		City:           foundUser.Address.City,
		Address:        foundUser.Address.StreetAddress,
		HouseNumber:    foundUser.Address.HouseNumber,
		Unit:           foundUser.Address.Unit,
		Status:         foundUser.PetSitter.Status,
		OnboardingStep: foundUser.PetSitter.OnboardingStep,
	}
	return r, nil

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
	CertificateKeys := make([]string, len(info.CertificateFiles))
	if info.CertificateFiles != nil {
		for i, file := range info.CertificateFiles {
			fileKeyValue := ps.getStorageKey(info.UserID)
			fileKey = &fileKeyValue
			if err := ps.storage.UploadFile(enums.PetSitterFile, *fileKey, file); err != nil {
				return err
			}
			CertificateKeys[i] = fileKeyValue
		}
	}
	FileKeys := make([]string, len(info.Files))
	if info.Files != nil {
		for i, file := range info.Files {
			fileKeyValue := ps.getStorageKey(info.UserID)
			fileKey = &fileKeyValue
			if err := ps.storage.UploadFile(enums.PetSitterFile, *fileKey, file); err != nil {
				return err
			}
			FileKeys[i] = *fileKey
		}
	}
	foundPetSitter.CertificateKeys = CertificateKeys
	foundPetSitter.FileKeys = FileKeys
	foundPetSitter.OnboardingStep = enums.OBS_Documents
	foundPetSitter.Status = enums.PSS_Draft
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.EditPetSitter(foundPetSitter)
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
	err = ps.CheckPetSitterStep(foundPetSitter.OnboardingStep, enums.OBS_Done)
	if err != nil {
		return nil, err
	}
	CertificateFiles := make([]string, len(foundPetSitter.CertificateKeys))
	for i, certKey := range foundPetSitter.CertificateKeys {
		certificateURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, certKey, 2)
		CertificateFiles[i] = certificateURL
	}
	Files := make([]string, len(foundPetSitter.FileKeys))
	for i, fileKey := range foundPetSitter.FileKeys {
		fileURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, fileKey, 2)
		Files[i] = fileURL
	}
	return &petsitter.DocumentResponse{
		CertificateFiles: CertificateFiles,
		Files:            Files,
	}, nil
}

func (ps *PetSitterService) SubmitSkills(SkillsInfo petsitter.SubmitSkillsRequest) error {
	_, err := ps.userService.FindVerifiedUserByID(SkillsInfo.UserID)
	if err != nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()

	foundPetSitter, err := ps.FindPetSitterByID(SkillsInfo.UserID)
	if err != nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = petSitterRepo.PreloadServices(foundPetSitter)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStep(foundPetSitter.OnboardingStep, enums.OBS_Done)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStatus(foundPetSitter.Status)
	if err != nil {
		return err
	}
	services := ps.GetPetsitterServicesResponse(SkillsInfo.Services, foundPetSitter.ID)
	foundPetSitter.Bio = &SkillsInfo.Bio
	foundPetSitter.Services = services
	foundPetSitter.PetKinds = append(foundPetSitter.PetKinds, SkillsInfo.PetKinds...)
	foundPetSitter.OnboardingStep = enums.OBS_Done
	foundPetSitter.Status = enums.PSS_InReview

	err = petSitterRepo.EditPetSitter(foundPetSitter)
	if err != nil {
		return err
	}
	return nil

}

func (ps *PetSitterService) GetPetsitterStatus(userID uint) (*petsitter.PetSitterStatusResponse, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundPetSitter, err := petSitterRepo.FindPetSitterByUserID(userID)
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
	foundPetSitter, err := petSitterRepo.FindPetSitterByUserID(id)
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

func (ps *PetSitterService) GetPetsitterServicesResponse(Services []enums.ServiceType, petSitterID uint) []entities.Service {
	r := make([]entities.Service, len(Services))
	for i, service := range Services {
		r[i] = entities.Service{
			PetSitterID: petSitterID,
			Type:        service,
		}
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

func (ps *PetSitterService) GetAllPetSitters(page, count int) (*petsitter.PetSittersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if count <= 0 || count > 100 {
		count = 10
	}

	offset := (page - 1) * count
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()

	total, err := petSitterRepo.GetPetSittersCount()
	if err != nil {
		return nil, err
	}

	petSitters, err := petSitterRepo.GetAllPetSitters(count, offset)
	if err != nil {
		return nil, err
	}

	userRepo := ps.unitOfWork.Factory().UserRepository()
	items := make([]petsitter.PetSitterListItemResponse, len(petSitters))

	for i, ps := range petSitters {
		user, err := userRepo.FindUserByID(ps.UserID)
		if err != nil {
			continue
		}
		items[i] = petsitter.PetSitterListItemResponse{
			ID:             ps.ID,
			UserID:         ps.UserID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          user.Email,
			PhoneNumber:    *user.Phone,
			Status:         ps.Status,
			OnboardingStep: ps.OnboardingStep,
			CreatedAt:      ps.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &petsitter.PetSittersListResponse{
		Total:      total,
		Page:       page,
		Count:      count,
		PetSitters: items,
	}, nil
}

func (ps *PetSitterService) SubmitPersonalInfo(petSitterInfo petsitter.SubmitPersonalInfoRequest) error {
	foundUser, err := ps.userService.FindVerifiedUserByID(petSitterInfo.UserID)
	if err != nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.User)
	}
	err = ps.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		userRepo := rf.UserRepository()
		addressRepo := rf.AddressRepository()
		err = userRepo.PreloadPetSitter(foundUser)
		if err != nil {
			return err
		}
		if foundUser.PetSitter == nil {
			return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
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
		addressInfo := address.AddressInfo{
			ProvinceName:  petSitterInfo.Province,
			CityName:      petSitterInfo.City,
			StreetAddress: petSitterInfo.Address,
			HouseNumber:   petSitterInfo.HouseNumber,
			Unit:          petSitterInfo.Unit,
		}
		createdAddress, err := ps.addressService.CreateAddressEntity(addressInfo)
		if err != nil {
			return err
		}
		if foundUser.Address != nil {
			foundAddress, err := addressRepo.FindUserAddressByUserID(foundUser.ID)
			if err != nil {
				return err
			}
			foundAddress.Province = createdAddress.Province
			foundAddress.City = createdAddress.City
			foundAddress.StreetAddress = createdAddress.StreetAddress
			foundAddress.HouseNumber = createdAddress.HouseNumber
			foundAddress.Unit = createdAddress.Unit
			err = addressRepo.Update(foundAddress)
			if err != nil {
				return err
			}
		} else {
			createdAddress.Refer = foundUser.ID
			createdAddress.Type = "User"
			err = addressRepo.Create(createdAddress)
			if err != nil {
				return err
			}
		}

		foundUser.Address = createdAddress
		foundUser.FirstName = petSitterInfo.FirstName
		foundUser.LastName = petSitterInfo.LastName
		foundUser.Email = petSitterInfo.Email
		foundUser.Gender = petSitterInfo.Gender
		foundUser.BirthDate = petSitterInfo.BirthDate
		foundUser.Phone = &petSitterInfo.Phone
		foundUser.PetSitter.Status = enums.PSS_Draft
		foundUser.PetSitter.OnboardingStep = enums.OBS_Profile
		err = userRepo.SaveUser(foundUser)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (ps *PetSitterService) GetCalendarSlotsResponse(calendarSlots []entities.CalendarSlot) []calendarslot.CalendarSlotInfoResponse {
	r := make([]calendarslot.CalendarSlotInfoResponse, len(calendarSlots))

	for i, slot := range calendarSlots {
		r[i] = calendarslot.CalendarSlotInfoResponse{
			ID:    slot.ID,
			Date:  slot.Date,
			Slots: slot.Slots,
		}
	}

	return r
}

func (ps *PetSitterService) GetFreeMap(calendarSlots []entities.CalendarSlot) map[string]map[interface{}]bool {
	availableSlots := make(map[string]map[interface{}]bool)

	for _, psSlot := range calendarSlots {
		if psSlot.Status == enums.Free {
			dateKey := psSlot.Date.Format("2006-01-02")

			if _, exists := availableSlots[dateKey]; !exists {
				availableSlots[dateKey] = make(map[interface{}]bool)
			}

			for _, slot := range psSlot.Slots {
				availableSlots[dateKey][slot] = true
			}
		}
	}

	return availableSlots
}

func (ps *PetSitterService) FindServiceByID(id uint) (*entities.Service, error) {
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	service, err := petSitterRepo.FindServiceByID(id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Service)
	}

	return service, nil
}

func (ps *PetSitterService) GetServiceResponse(serviceEntity *entities.Service) servicedto.ServiceInfoResponse {
	return servicedto.ServiceInfoResponse{
		ID:          serviceEntity.ID,
		Type:        serviceEntity.Type.String(),
		Description: serviceEntity.Description,
		Price:       serviceEntity.Price,
	}
}
