package service

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/address"
	calendarslot "hona/backend/internal/application/dto/calendar_slot"
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/petsitter"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainstorage "hona/backend/internal/domain/storage"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
	"sort"
	"time"

	"github.com/lib/pq"
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

func (ps *PetSitterService) GetCalendarSlots(info calendarslot.GetCalendarSlotsRequest) ([]calendarslot.CalendarSlotInfoResponse, error) {
	petSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return nil, err
	}
	if err := ps.PreloadFields(petSitter, []string{"Schedule"}); err != nil {
		return nil, err
	}
	if petSitter.Schedule == nil {
		return []calendarslot.CalendarSlotInfoResponse{}, nil
	}

	start := time.Now()
	end := start.AddDate(0, 0, 30)
	filtered := make([]entities.CalendarSlot, 0, len(petSitter.Schedule))
	for _, slot := range petSitter.Schedule {
		if slot.Date.Before(start) || slot.Date.After(end) {
			continue
		}
		filtered = append(filtered, slot)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Date.Before(filtered[j].Date)
	})

	return ps.GetCalendarSlotsResponse(filtered), nil
}

func (ps *PetSitterService) UpdateFreeCalendarSlots(info calendarslot.UpdateFreeCalendarSlotsRequest) error {
	if len(info.Add) == 0 && len(info.Remove) == 0 {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ve
	}

	petSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return err
	}
	if err := ps.PreloadFields(petSitter, []string{"Schedule"}); err != nil {
		return err
	}

	addSlots := ps.makeFreeCalendarSlots(info.Add)
	removeSlots := ps.makeFreeCalendarSlots(info.Remove)
	updatedSchedule, err := applyFreeSlotPatch(petSitter.Schedule, addSlots, removeSlots)
	if err != nil {
		return err
	}

	petSitter.Schedule = updatedSchedule
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.ReplaceSchedule(petSitter, petSitter.Schedule)
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

func (ps *PetSitterService) UpdateRatingAndCommentsCount(petSitterID uint, rating float32, commentsCount uint) error {
	foundPetSitter, err := ps.GetPetSitterByID(petSitterID)
	if err != nil {
		return err
	}

	foundPetSitter.Rating = rating
	foundPetSitter.CommentsCount = commentsCount

	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.UpdatePetSitter(foundPetSitter)
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
	return petSitterRepo.ReplaceSchedule(petSitter, petSitter.Schedule)

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
		return nil, exceptions.NewNotVerifiedError()
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
		return nil, exceptions.NewNotVerifiedError()
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = userRepo.PreloadAddress(foundUser)
	if err != nil && foundUser.Address != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Address)
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Profile)
	if err != nil {
		return nil, err
	}
	type Address struct {
		Province      enums.Province
		City          enums.City
		StreetAddress string
		HouseNumber   uint
		Unit          uint
	}
	address := &Address{}
	if foundUser.Address == nil {
		address.Province = 0
		address.City = 0
		address.StreetAddress = ""
		address.HouseNumber = 0
		address.Unit = 0
	}
	if foundUser.Address != nil {
		address.Province = foundUser.Address.Province
		address.City = foundUser.Address.City
		address.StreetAddress = foundUser.Address.StreetAddress
		address.HouseNumber = foundUser.Address.HouseNumber
		address.Unit = foundUser.Address.Unit
	}
	phoneNumber := ""
	if foundUser.Phone != nil {
		phoneNumber = *foundUser.Phone
	}
	var userProfilePicURL *string
	if foundUser.PictureLink != nil {
		url, _ := ps.storage.GetPresignedURL(enums.UserProfilePic, *foundUser.PictureLink, 2)
		userProfilePicURL = &url
	}
	r := &petsitter.PersonalInfoResponse{
		FirstName:      foundUser.FirstName,
		LastName:       foundUser.LastName,
		Email:          foundUser.Email,
		PhoneNumber:    phoneNumber,
		BirthDate:      foundUser.BirthDate,
		Gender:         foundUser.Gender,
		Province:       address.Province,
		City:           address.City,
		Address:        address.StreetAddress,
		HouseNumber:    address.HouseNumber,
		Unit:           address.Unit,
		UserProfilePic: userProfilePicURL,
		Status:         foundUser.PetSitter.Status,
		OnboardingStep: foundUser.PetSitter.OnboardingStep,
	}
	return r, nil

}

func (ps *PetSitterService) UploadDocuments(info petsitter.UploadDocumentsRequest) error {
	foundUser, err := ps.userService.FindVerifiedUserByID(info.UserID)
	if err != nil {
		return exceptions.NewNotVerifiedError()
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Documents)
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
	foundUser.PetSitter.CertificateKeys = CertificateKeys
	foundUser.PetSitter.FileKeys = FileKeys
	foundUser.PetSitter.OnboardingStep = enums.OBS_Documents
	foundUser.PetSitter.Status = enums.PSS_Draft
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.UpdatePetSitter(foundUser.PetSitter)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PetSitterService) GetDocuments(userID uint) (*petsitter.DocumentResponse, error) {
	foundUser, err := ps.userService.FindVerifiedUserByID(userID)
	if err != nil {
		return nil, exceptions.NewNotVerifiedError()
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return nil, err
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Done)
	if err != nil {
		return nil, err
	}
	CertificateFiles := make([]string, len(foundUser.PetSitter.CertificateKeys))
	for i, certKey := range foundUser.PetSitter.CertificateKeys {
		certificateURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, certKey, 2)
		CertificateFiles[i] = certificateURL
	}
	Files := make([]string, len(foundUser.PetSitter.FileKeys))
	for i, fileKey := range foundUser.PetSitter.FileKeys {
		fileURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, fileKey, 2)
		Files[i] = fileURL
	}
	return &petsitter.DocumentResponse{
		CertificateFiles: CertificateFiles,
		Files:            Files,
	}, nil
}

func (ps *PetSitterService) SubmitSkills(SkillsInfo petsitter.SubmitSkillsRequest) error {
	foundUser, err := ps.userService.FindVerifiedUserByID(SkillsInfo.UserID)
	if err != nil {
		return exceptions.NewNotVerifiedError()
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	err = ps.CheckPetSitterStep(foundUser.PetSitter.OnboardingStep, enums.OBS_Done)
	if err != nil {
		return err
	}
	err = ps.CheckPetSitterStatus(foundUser.PetSitter.Status)
	if err != nil {
		return err
	}

	services := ps.GetPetsitterServicesResponse(SkillsInfo.Services, foundUser.PetSitter.ID)
	foundUser.PetSitter.Bio = &SkillsInfo.Bio
	foundUser.PetSitter.Services = services
	foundUser.PetSitter.PetKinds = SkillsInfo.PetKinds
	foundUser.PetSitter.OnboardingStep = enums.OBS_Done
	foundUser.PetSitter.Status = enums.PSS_InReview
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.UpdatePetSitter(foundUser.PetSitter)

	if err != nil {
		return err
	}
	return nil

}
func mapPetKinds(kinds []enums.PetKind) pq.Int32Array {
	result := make(pq.Int32Array, len(kinds))
	for i, k := range kinds {
		result[i] = int32(k)
	}
	return result
}

func (ps *PetSitterService) GetPetsitterStatus(userID uint) (*petsitter.PetSitterStatusResponse, error) {
	foundUser, err := ps.userService.FindVerifiedUserByID(userID)
	if err != nil {
		return nil, exceptions.NewNotVerifiedError()
	}
	userRepo := ps.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(foundUser)
	if err != nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return &petsitter.PetSitterStatusResponse{
		OnboardingStep: foundUser.PetSitter.OnboardingStep,
		Status:         foundUser.PetSitter.Status,
	}, nil
}

func (ps *PetSitterService) GetPetKinds(info petsitter.GetPetKindsRequest) ([]pet.PetKindResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return nil, err
	}
	return ps.buildPetKindsResponse([]enums.PetKind(foundPetSitter.PetKinds)), nil
}

func (ps *PetSitterService) UpdatePetKinds(info petsitter.UpdatePetKindsRequest) error {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return err
	}
	foundPetSitter.PetKinds = entities.PetKinds(info.PetKinds)
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.UpdatePetSitter(foundPetSitter)
}

func (ps *PetSitterService) GetServices(info petsitter.GetServicesRequest) ([]servicedto.ServiceInfoResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return nil, err
	}
	if err := ps.PreloadFields(foundPetSitter, []string{"Services"}); err != nil {
		return nil, err
	}
	return ps.GetServicesResponse(foundPetSitter)
}

func (ps *PetSitterService) CreateService(info petsitter.CreateServiceRequest) (servicedto.ServiceInfoResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return servicedto.ServiceInfoResponse{}, err
	}
	service := entities.Service{
		PetSitterID: foundPetSitter.ID,
		Type:        info.Type,
		Price:       info.Price,
		Description: info.Description,
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	if err := petSitterRepo.CreateService(&service); err != nil {
		return servicedto.ServiceInfoResponse{}, err
	}
	return ps.GetServiceResponse(&service), nil
}

func (ps *PetSitterService) UpdateService(info petsitter.UpdateServiceRequest) (servicedto.ServiceInfoResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return servicedto.ServiceInfoResponse{}, err
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundService, err := petSitterRepo.FindServiceByID(info.ID)
	if err != nil {
		return servicedto.ServiceInfoResponse{}, err
	}
	if foundService == nil {
		return servicedto.ServiceInfoResponse{}, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Service)
	}
	if foundService.PetSitterID != foundPetSitter.ID {
		return servicedto.ServiceInfoResponse{}, exceptions.NewAccessDeniedError("can't update other's services")
	}
	foundService.Type = info.Type
	foundService.Price = info.Price
	foundService.Description = info.Description
	if err := petSitterRepo.UpdateService(foundService); err != nil {
		return servicedto.ServiceInfoResponse{}, err
	}
	return ps.GetServiceResponse(foundService), nil
}

func (ps *PetSitterService) DeleteService(info petsitter.DeleteServiceRequest) error {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.UserID)
	if err != nil {
		return err
	}
	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	foundService, err := petSitterRepo.FindServiceByID(info.ID)
	if err != nil {
		return err
	}
	if foundService == nil {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Service)
	}
	if foundService.PetSitterID != foundPetSitter.ID {
		return exceptions.NewAccessDeniedError("can't delete other's services")
	}
	return petSitterRepo.DeleteService(foundService)
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
	r := make([]entities.Service, 0)
	for _, service := range Services {
		r = append(r, entities.Service{
			PetSitterID: petSitterID,
			Type:        service,
			Kind:        "PetSitter",
		})
	}
	return r
}

func (ps *PetSitterService) CheckPetSitterStatus(pss enums.PetSitterStatus) error {
	if pss == enums.PSS_Rejected {
		return exceptions.NewForbiddenError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus, bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if pss == enums.PSS_Suspended {
		return exceptions.NewForbiddenError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus, bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if pss == enums.PSS_Active {
		return exceptions.NewForbiddenError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus, bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	if pss == enums.PSS_InReview {
		return exceptions.NewForbiddenError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus, bootstrap.Run().Constants.ErrorFields.PetSitter)
	}
	return nil
}

func (ps *PetSitterService) CheckPetSitterStep(currentStep enums.OnboardingStep, requiredStep enums.OnboardingStep) error {
	if (currentStep + 1) < requiredStep {
		return exceptions.NewForbiddenError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus, bootstrap.Run().Constants.ErrorFields.PetSitter)
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
		phoneNumber := ""
		if user.Phone != nil {
			phoneNumber = *user.Phone
		}
		items[i] = petsitter.PetSitterListItemResponse{
			ID:             ps.ID,
			UserID:         ps.UserID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          user.Email,
			PhoneNumber:    phoneNumber,
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
		return exceptions.NewNotVerifiedError()
	}
	errr := ps.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
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

		var ProfilePicKey *string
		if petSitterInfo.UserProfilePic != nil {
			profileKeyValue := ps.getStorageKey(foundUser.ID)
			ProfilePicKey = &profileKeyValue
			if err := ps.storage.UploadFile(enums.UserProfilePic, *ProfilePicKey, petSitterInfo.UserProfilePic); err != nil {
				return err
			}
			foundUser.PictureLink = ProfilePicKey
		}

		err = userRepo.SaveUser(foundUser)
		if err != nil {
			return err
		}

		// Save PetSitter changes explicitly
		foundUser.PetSitter.Status = enums.PSS_Draft
		foundUser.PetSitter.OnboardingStep = enums.OBS_Profile
		petSitterRepo := rf.PetSitterRepository()
		err = petSitterRepo.UpdatePetSitter(foundUser.PetSitter)
		if err != nil {
			return err
		}
		return nil
	})
	return errr
}

func (ps *PetSitterService) GetCalendarSlotsResponse(calendarSlots []entities.CalendarSlot) []calendarslot.CalendarSlotInfoResponse {
	r := make([]calendarslot.CalendarSlotInfoResponse, len(calendarSlots))

	for i, slot := range calendarSlots {
		r[i] = calendarslot.CalendarSlotInfoResponse{
			ID:     slot.ID,
			Date:   slot.Date,
			Slots:  []enums.Slot(slot.Slots),
			Status: slot.Status,
		}
	}

	return r
}

func (ps *PetSitterService) makeFreeCalendarSlots(slots []calendarslot.CalendarSlotRequest) []entities.CalendarSlot {
	result := make([]entities.CalendarSlot, len(slots))
	for i, slot := range slots {
		result[i] = entities.CalendarSlot{
			Date:   slot.Date,
			Slots:  entities.Slots(slot.Slots),
			Status: enums.Free,
		}
	}
	return result
}

func applyFreeSlotPatch(schedule []entities.CalendarSlot, addSlots []entities.CalendarSlot, removeSlots []entities.CalendarSlot) ([]entities.CalendarSlot, error) {
	freeSlots, bookedSlots, dateByKey := buildSlotMaps(schedule)

	if err := applyFreeSlotAdds(addSlots, freeSlots, bookedSlots, dateByKey); err != nil {
		return nil, err
	}
	if err := applyFreeSlotRemovals(removeSlots, freeSlots, bookedSlots, dateByKey); err != nil {
		return nil, err
	}

	return rebuildSchedule(schedule, freeSlots, dateByKey), nil
}

func buildSlotMaps(schedule []entities.CalendarSlot) (map[string]map[enums.Slot]bool, map[string]map[enums.Slot]bool, map[string]time.Time) {
	freeSlots := make(map[string]map[enums.Slot]bool)
	bookedSlots := make(map[string]map[enums.Slot]bool)
	dateByKey := make(map[string]time.Time)

	for _, slot := range schedule {
		dateKey := slot.Date.Format("2006-01-02")
		if _, exists := dateByKey[dateKey]; !exists {
			dateByKey[dateKey] = slot.Date
		}
		switch slot.Status {
		case enums.Booked:
			addSlotsToMap(bookedSlots, dateKey, []enums.Slot(slot.Slots))
		case enums.Free:
			addSlotsToMap(freeSlots, dateKey, []enums.Slot(slot.Slots))
		}
	}

	return freeSlots, bookedSlots, dateByKey
}

func addSlotsToMap(target map[string]map[enums.Slot]bool, dateKey string, slots []enums.Slot) {
	if _, exists := target[dateKey]; !exists {
		target[dateKey] = make(map[enums.Slot]bool)
	}
	for _, s := range slots {
		target[dateKey][s] = true
	}
}

func applyFreeSlotAdds(addSlots []entities.CalendarSlot, freeSlots map[string]map[enums.Slot]bool, bookedSlots map[string]map[enums.Slot]bool, dateByKey map[string]time.Time) error {
	for _, slot := range addSlots {
		dateKey := slot.Date.Format("2006-01-02")
		dateByKey[dateKey] = slot.Date
		if hasBookedConflict(bookedSlots, dateKey, []enums.Slot(slot.Slots)) {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.CalendarConflict)
			return &ce
		}
		addSlotsToMap(freeSlots, dateKey, []enums.Slot(slot.Slots))
	}
	return nil
}

func applyFreeSlotRemovals(removeSlots []entities.CalendarSlot, freeSlots map[string]map[enums.Slot]bool, bookedSlots map[string]map[enums.Slot]bool, dateByKey map[string]time.Time) error {
	for _, slot := range removeSlots {
		dateKey := slot.Date.Format("2006-01-02")
		dateByKey[dateKey] = slot.Date
		if hasBookedConflict(bookedSlots, dateKey, []enums.Slot(slot.Slots)) {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.CalendarConflict)
			return &ce
		}
		daySlots, exists := freeSlots[dateKey]
		for _, s := range []enums.Slot(slot.Slots) {
			if !exists || !daySlots[s] {
				var ve exceptions.ValidationErrors
				ve.AddError(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
				return &ve
			}
			delete(daySlots, s)
		}
		if len(daySlots) == 0 {
			delete(freeSlots, dateKey)
		}
	}
	return nil
}

func hasBookedConflict(bookedSlots map[string]map[enums.Slot]bool, dateKey string, slots []enums.Slot) bool {
	day, exists := bookedSlots[dateKey]
	if !exists {
		return false
	}
	for _, s := range slots {
		if day[s] {
			return true
		}
	}
	return false
}

func rebuildSchedule(schedule []entities.CalendarSlot, freeSlots map[string]map[enums.Slot]bool, dateByKey map[string]time.Time) []entities.CalendarSlot {
	updated := make([]entities.CalendarSlot, 0, len(schedule))
	for _, slot := range schedule {
		if slot.Status == enums.Booked {
			updated = append(updated, slot)
		}
	}

	for dateKey, slotsMap := range freeSlots {
		if len(slotsMap) == 0 {
			continue
		}
		slots := make([]enums.Slot, 0, len(slotsMap))
		for s := range slotsMap {
			slots = append(slots, s)
		}
		sort.Slice(slots, func(i, j int) bool {
			return slots[i] < slots[j]
		})
		updated = append(updated, entities.CalendarSlot{
			Date:   dateByKey[dateKey],
			Slots:  entities.Slots(slots),
			Status: enums.Free,
		})
	}

	sort.Slice(updated, func(i, j int) bool {
		return updated[i].Date.Before(updated[j].Date)
	})

	return updated
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
func extractServiceType(filters []general.Filter) (*enums.ServiceType, error) {
	for _, f := range filters {
		if f.Field != "serviceType" {
			continue
		}
		// value may come as float64 (json), int, etc.
		switch v := f.Value.(type) {
		case float64:
			st := enums.ServiceType(uint(v))
			return &st, nil
		case int:
			st := enums.ServiceType(uint(v))
			return &st, nil
		case int32:
			st := enums.ServiceType(uint(v))
			return &st, nil
		case int64:
			st := enums.ServiceType(uint(v))
			return &st, nil
		case uint:
			st := enums.ServiceType(v)
			return &st, nil
		case uint64:
			st := enums.ServiceType(uint(v))
			return &st, nil
		default:
			return nil, fmt.Errorf("invalid serviceType value type: %T", f.Value)
		}
	}
	return nil, nil
}
func extractMinPrice(filters []general.Filter) *uint {
	if filters == nil {
		return nil
	}

	for _, f := range filters {
		if f.Field == "minPrice" {
			if v, ok := toUint(f.Value); ok {
				return &v
			}
		}
	}
	return nil
}
func extractMaxPrice(filters []general.Filter) *uint {
	if filters == nil {
		return nil
	}

	for _, f := range filters {
		if f.Field == "maxPrice" {
			if v, ok := toUint(f.Value); ok {
				return &v
			}
		}
	}
	return nil
}
func toUint(value interface{}) (uint, bool) {
	switch v := value.(type) {
	case float64:
		return uint(v), true
	case int:
		return uint(v), true
	case int32:
		return uint(v), true
	case int64:
		return uint(v), true
	case uint:
		return v, true
	case uint64:
		return uint(v), true
	default:
		return 0, false
	}
}
func (ps *PetSitterService) SearchPetSitters(info petsitter.SearchPetSittersRequest) ([]*petsitter.SearchPetSitterInfoResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(info.Limit, info.Offset).
		WithSorting(info.Sorts).
		WithFilters(info.Filters)

	serviceType, err := extractServiceType(info.Filters)
	if err != nil {
		return nil, 0, err
	}

	minFilter := extractMinPrice(info.Filters)
	maxFilter := extractMaxPrice(info.Filters)

	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	petSitters, total, err := petSitterRepo.SearchPetSitters(options)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*petsitter.SearchPetSitterInfoResponse, 0, len(petSitters))

	for _, petSitter := range petSitters {
		user, err := ps.userService.FindUserByID(petSitter.UserID)
		if err != nil {
			return nil, 0, err
		}
		if err := ps.userService.PreloadFields(user, []string{"Address"}); err != nil {
			return nil, 0, err
		}

		var province, city string
		if user.Address != nil {
			province = user.Address.Province.String()
			city = user.Address.City.String()
		}

		if err := petSitterRepo.PreloadFields(petSitter, []string{"Comments"}); err != nil {
			return nil, 0, err
		}

		var comments uint
		var rateSum uint
		for _, c := range petSitter.Comments {
			comments++
			rateSum += c.Rating
		}
		var averageRating float64
		if comments > 0 {
			averageRating = float64(rateSum) / float64(comments)
		}

		if err := petSitterRepo.PreloadServices(petSitter); err != nil {
			return nil, 0, err
		}

		services := make([]string, len(petSitter.Services))
		for i, s := range petSitter.Services {
			services[i] = s.Type.String()
		}


		var minPrice uint
		if len(petSitter.Services) > 0 {
			minPrice = 0
			for _, s := range petSitter.Services {
				if s.Kind != "petSitter" || s.DeletedAt.Valid {
					continue
				}
				if serviceType != nil && s.Type != *serviceType {
					continue
				}
				if minFilter != nil && s.Price < *minFilter {
					continue
				}
				if maxFilter != nil && s.Price > *maxFilter {
					continue
				}
				if minPrice == 0 || s.Price < minPrice {
					minPrice = s.Price
				}
			}
		}

		petkinds := make([]string, len(petSitter.PetKinds))
		for j, pk := range petSitter.PetKinds {
			petkinds[j] = pk.String()
		}

		res = append(res, &petsitter.SearchPetSitterInfoResponse{
			ID:        petSitter.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Province:  province,
			City:      city,
			Services:  services,
			PetKinds:  petkinds,
			MinPrice:  minPrice,
			Rate:      averageRating,
			Comments:  comments,
		})
	}

	return res, total, nil
}

func (ps *PetSitterService) SearchPetSittersForAdmin(info petsitter.AdminSearchPetSittersRequest) ([]petsitter.PetSitterListItemResponse, int64, error) {
	options := postgres.NewQueryOptions().WithPagination(info.Limit, info.Offset)
	if len(info.Filters) > 0 {
		options.WithFilters(info.Filters)
	}
	if len(info.Sorts) > 0 {
		options.WithSorting(info.Sorts)
	}

	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	petSitters, total, err := petSitterRepo.SearchPetSitters(options)
	if err != nil {
		return nil, 0, err
	}

	items := make([]petsitter.PetSitterListItemResponse, len(petSitters))
	for i, psr := range petSitters {
		user, err := ps.userService.FindUserByID(psr.UserID)
		if err != nil {
			continue
		}
		phoneNumber := ""
		if user.Phone != nil {
			phoneNumber = *user.Phone
		}
		items[i] = petsitter.PetSitterListItemResponse{
			ID:             psr.ID,
			UserID:         psr.UserID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Email:          user.Email,
			PhoneNumber:    phoneNumber,
			Status:         psr.Status,
			OnboardingStep: psr.OnboardingStep,
			CreatedAt:      psr.CreatedAt.String(),
		}
	}

	return items, total, nil
}

func (ps *PetSitterService) GetPetSitterDetails(info petsitter.GetPetSitterDetailsRequest) (*petsitter.PetSitterDetailsResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}
	err = ps.PreloadFields(foundPetSitter, []string{"Services"})
	if err != nil {
		return nil, err
	}

	foundUser, err := ps.userService.FindUserByID(foundPetSitter.UserID)
	if err != nil {
		return nil, err
	}

	personalInfo, err := ps.userService.GetUserInfoResponse(foundUser)
	if err != nil {
		return nil, err
	}
	services, err := ps.GetServicesResponse(foundPetSitter)
	if err != nil {
		return nil, err
	}

	petKinds := ps.buildPetKindsResponse([]enums.PetKind(foundPetSitter.PetKinds))
	documents := ps.GetDocumentsInfo(foundPetSitter)

	return &petsitter.PetSitterDetailsResponse{
		PersonalInfo: *personalInfo,
		Skills: petsitter.SkillsResponse{
			Bio:      *foundPetSitter.Bio,
			Services: services,
			PetKinds: petKinds,
		},
		Documents:      documents,
		Status:         foundPetSitter.Status,
		OnboardingStep: foundPetSitter.OnboardingStep,
		CreatedAt:      foundPetSitter.CreatedAt.String(),
	}, nil
}

func (ps *PetSitterService) GetPetSitterProfile(info petsitter.GetPetSitterProfileRequest) (*petsitter.PetSitterProfileResponse, error) {
	foundPetSitter, err := ps.GetPetSitterByID(info.PetSitterID)
	if err != nil {
		return nil, err
	}
	if err := ps.PreloadFields(foundPetSitter, []string{"Services"}); err != nil {
		return nil, err
	}

	foundUser, err := ps.userService.FindUserByID(foundPetSitter.UserID)
	if err != nil {
		return nil, err
	}
	if err := ps.userService.PreloadFields(foundUser, []string{"Address"}); err != nil {
		return nil, err
	}

	services, err := ps.GetServicesResponse(foundPetSitter)
	if err != nil {
		return nil, err
	}

	province := ""
	city := ""
	if foundUser.Address != nil {
		province = foundUser.Address.Province.String()
		city = foundUser.Address.City.String()
	}

	bio := ""
	if foundPetSitter.Bio != nil {
		bio = *foundPetSitter.Bio
	}

	petKinds := ps.buildPetKindsResponse([]enums.PetKind(foundPetSitter.PetKinds))

	return &petsitter.PetSitterProfileResponse{
		ID:          foundPetSitter.ID,
		UserID:      foundPetSitter.UserID,
		FirstName:   foundUser.FirstName,
		LastName:    foundUser.LastName,
		PictureLink: foundUser.PictureLink,
		Province:    province,
		City:        city,
		Bio:         bio,
		Services:    services,
		PetKinds:    petKinds,
		CreatedAt:   foundPetSitter.CreatedAt.String(),
	}, nil
}

func (ps *PetSitterService) buildPetKindsResponse(petKinds []enums.PetKind) []pet.PetKindResponse {
	res := make([]pet.PetKindResponse, len(petKinds))
	for i, pk := range petKinds {
		res[i] = pet.PetKindResponse{
			Num:  pk,
			Name: pk.String(),
		}
	}
	return res
}

func (ps *PetSitterService) ChangePetSitterStatus(info petsitter.ChangePetSitterStatusRequest) error {
	foundPetSitter, err := ps.GetPetSitterByUserID(info.PetSitterUserID)
	if err != nil {
		return err
	}
	switch info.Status {
	case enums.PSS_Suspended:
		if foundPetSitter.Status != enums.PSS_Active && foundPetSitter.Status != enums.PSS_Rejected {
			return exceptions.NewAccessDeniedError("can't change the status to suspend")
		}
	case enums.PSS_Active, enums.PSS_Rejected:
		if foundPetSitter.Status != enums.PSS_Suspended && foundPetSitter.Status != enums.PSS_InReview {
			return exceptions.NewAccessDeniedError("can't change the status to active or rejected")
		}
	}

	foundPetSitter.Status = info.Status

	petSitterRepo := ps.unitOfWork.Factory().PetSitterRepository()
	return petSitterRepo.UpdatePetSitter(foundPetSitter)
}

func (ps *PetSitterService) GetDocumentsInfo(petSitter *entities.PetSitter) petsitter.DocumentResponse {
	var CertificateFiles []string
	if petSitter.CertificateKeys != nil {
		CertificateFiles = make([]string, len(petSitter.CertificateKeys))
		for i, certKey := range petSitter.CertificateKeys {
			certificateURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, certKey, 2)
			CertificateFiles[i] = certificateURL
		}
	}
	var Files []string
	if petSitter.FileKeys != nil {
		Files = make([]string, len(petSitter.FileKeys))
		for i, fileKey := range petSitter.FileKeys {
			fileURL, _ := ps.storage.GetPresignedURL(enums.PetSitterFile, fileKey, 2)
			Files[i] = fileURL
		}
	}
	return petsitter.DocumentResponse{
		CertificateFiles: CertificateFiles,
		Files:            Files,
	}
}
