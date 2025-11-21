package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/request"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	"sort"
	"time"
)

type RequestService struct {
	userService         usecase.UserService
	provinceService     usecase.ProvinceService
	addressService      usecase.AddressService
	petService          usecase.PetService
	serviceService      usecase.ServiceService
	petSitterService    usecase.PetSitterService
	calendarSlotService usecase.CalendarSlotService
	unitOfWork          ports.UnitOfWork
}

func NewRequestService(userService usecase.UserService, unitOfWork ports.UnitOfWork, provinceService usecase.ProvinceService, addressService usecase.AddressService, petService usecase.PetService, calendarSlotService usecase.CalendarSlotService) *RequestService {
	return &RequestService{
		userService:         userService,
		unitOfWork:          unitOfWork,
		provinceService:     provinceService,
		addressService:      addressService,
		petService:          petService,
		calendarSlotService: calendarSlotService,
	}
}

// TODO: rethink error handlings

func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	petSitter, err := rs.petSitterService.GetPetSitterByID(info.PetSitterUserID)
	if err != nil {
		return err
	}

	err = rs.petSitterService.PreloadFields(petSitter, []string{"Schedule", "Services"})
	if err != nil {
		return err
	}

	calendarSlots := rs.makeCalendarSlots(info.CalenderSlots)
	err = rs.validateCalendarSlots(petSitter.Schedule, calendarSlots)
	if err != nil {
		return err
	}

	user, err := rs.userService.FindVerifiedUserByID(info.UserID)
	if err != nil {
		return err
	}

	err = rs.userService.PreloadFields(user, []string{"Pets"})
	if err != nil {
		return err
	}

	pets, err := rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	if info.AddressInfo != nil {
		madeAddress, err := rs.addressService.CreateAddress(*info.AddressInfo)
		if err != nil {
			return err
		}
		address = madeAddress
	} else if info.AddressID != nil {
		address, err = rs.addressService.FindAddressByID(*info.AddressID)
		if err != nil {
			return err
		}
	} else {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Address, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ve
	}

	serviceEntity, err := rs.petSitterService.ValidateService(petSitter.Services, info.ServiceID)
	if err != nil {
		return err
	}

	totalPrice := rs.calculateTotalPrice(serviceEntity, petSitter.Schedule, len(pets))

	newRequest := &entities.Request{
		UserID:          info.UserID,
		PetSitterUserID: info.PetSitterUserID,
		Status:          enums.Pending,
		Chat:            entities.Chat{},
		TransferID:      nil,
		CalendarSlots:   calendarSlots,
		Pets:            pets,
		TotalPrice:      totalPrice,
		Notes:           info.Notes,
		Comment:         nil,
		Address:         *address,
		Service:         *serviceEntity,
	}

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	err = requestRepo.CreateRequest(newRequest)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RequestService) GetCreateRequestInfo(info request.GetCreateRequestInfoRequest) (*request.CreateRequestInfoResponse, error) {
	addresses, err := rs.addressService.GetUserAddressesInfo(info.UserID)
	if err != nil {
		return nil, err
	}

	foundUser, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return nil, err
	}

	err = rs.userService.PreloadFields(foundUser, []string{"Pets"})
	if err != nil {
		return nil, err
	}

	petsData, err := rs.petService.GetPetsBasicDataResponse(foundUser.Pets)
	if err != nil {
		return nil, err
	}

	petSitter, err := rs.petSitterService.GetPetSitterByID(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	err = rs.petSitterService.PreloadFields(petSitter, []string{"Schedule", "Services"})
	if err != nil {
		return nil, err
	}

	freeSlots, err := rs.petSitterService.GetPetSitterFreeSlotsResponse(petSitter)
	if err != nil {
		return nil, err
	}

	servicesData, err := rs.petSitterService.GetAvailableServicesResponse(petSitter)
	if err != nil {
		return nil, err
	}

	return &request.CreateRequestInfoResponse{
		Services:          servicesData,
		Addresses:         addresses,
		Pets:              petsData,
		FreeCalendarSlots: freeSlots,
	}, nil
}

func (rs *RequestService) EditRequest(info request.EditRequestRequest) error {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return err
	}

	if foundRequest.Status != enums.Pending {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterUserID)
	if err != nil {
		return err
	}

	err = rs.petSitterService.PreloadFields(petSitter, []string{"Schedule", "Services"})
	if err != nil {
		return err
	}

	user, err := rs.userService.FindVerifiedUserByID(info.UserID)
	if err != nil {
		return err
	}

	pets, err := rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	if info.AddressInfo != nil {
		madeAddress, err := rs.addressService.CreateAddress(*info.AddressInfo)
		if err != nil {
			return err
		}
		address = madeAddress
	} else if info.AddressID != nil {
		address, err = rs.addressService.FindAddressByID(*info.AddressID)
		if err != nil {
			return err
		}
	} else {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Address, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ve
	}

	serviceEntity, err := rs.petSitterService.ValidateService(petSitter.Services, info.ServiceID)
	if err != nil {
		return err
	}

	totalPrice := rs.calculateTotalPrice(serviceEntity, petSitter.Schedule, len(pets))

	calendarSlots := rs.makeCalendarSlots(info.CalenderSlots)

	foundRequest.CalendarSlots = calendarSlots
	foundRequest.Pets = pets
	foundRequest.Notes = info.Notes
	foundRequest.Address = *address
	foundRequest.Service = *serviceEntity
	foundRequest.TotalPrice = uint(totalPrice)

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	requestRepo.EditRequest(foundRequest)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RequestService) CancelRequest(info request.CancelRequestRequest) error {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return nil
	}

	err = rs.PreloadFields(foundRequest, []string{"CalendarSlots"})
	if err != nil {
		return err
	}

	if foundRequest.UserID != info.UserID || foundRequest.PetSitterUserID == info.UserID {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	if foundRequest.Status == enums.Finished || foundRequest.Status == enums.Canceled {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	sort.Slice(foundRequest.CalendarSlots, func(i, j int) bool {
		return foundRequest.CalendarSlots[i].Date.Before(foundRequest.CalendarSlots[j].Date)
	})

	if foundRequest.CalendarSlots[0].Date.Before(time.Now().Add(-time.Hour * 24)) {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	foundRequest.Status = enums.Canceled

	// TODO: change sitter slots
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	requestRepo.EditRequest(foundRequest)

	return nil
}

func (rs *RequestService) GetRequestFullData(info request.GetRequestFullDataRequest) (*request.RequestFullDataResponse, error) {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return nil, err
	}

	err = rs.PreloadFields(foundRequest, []string{"Pets", "CalendarSlots", "Service", "Address"})
	if err != nil {
		return nil, err
	}

	petsData, err := rs.petService.GetPetsBasicDataResponse(foundRequest.Pets)
	if err != nil {
		return nil, err
	}

	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	err = rs.petSitterService.PreloadFields(petSitter, []string{"Schedule"})
	if err != nil {
		return nil, err
	}

	err = rs.validateCalendarSlots(petSitter.Schedule, foundRequest.CalendarSlots)
	if err != nil {
		foundRequest.Status = enums.Conflict
		requestRepo := rs.unitOfWork.Factory().RequestRepository()
		err = requestRepo.EditRequest(foundRequest)
		if err != nil {
			return nil, err
		}
	}

	return &request.RequestFullDataResponse{
		RequestID:       foundRequest.ID,
		PetSitterUserID: foundRequest.PetSitterUserID,
		Service:         rs.serviceService.GetServiceResponse(&foundRequest.Service),
		Pets:            petsData,
		Address:         rs.addressService.GetUserAddressInfo(&foundRequest.Address),
		Notes:           foundRequest.Notes,
		TotalPrice:      foundRequest.TotalPrice,
		Status:          foundRequest.Status.String(),
		TransferID:      foundRequest.TransferID,
		CalendarSlots:   rs.calendarSlotService.GetCalendarSlotsResponse(foundRequest.CalendarSlots),
	}, nil
}

func (rs *RequestService) RespondToRequest(info request.RespondToRequestRequest) error {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return err
	}
	err = rs.PreloadFields(foundRequest, []string{"CalendarSlots"})
	if err != nil {
		return err
	}
	if foundRequest.Status != enums.Pending {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}
	if foundRequest.PetSitterUserID != info.UserID {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	if info.Accept {
		foundRequest.Status = enums.Accepted

		petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterUserID)
		if err != nil {
			return err
		}
		err = rs.validateCalendarSlots(petSitter.Schedule, foundRequest.CalendarSlots)
		if err != nil {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.AlreadyExist)
		}
	} else {
		foundRequest.Status = enums.Dismissed
	}

	// TODO: Update PetSitter Slots if accepted
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	requestRepo.EditRequest(foundRequest)

	return nil
}

func (rs *RequestService) validateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []entities.CalendarSlot) error {
	availableSlots := rs.calendarSlotService.GetFreeMap(petSitterSlots)

	for _, req := range requestSlots {
		dateKey := req.Date.Format("2006-01-02")

		daySlots, dayExists := availableSlots[dateKey]

		for _, userSlot := range req.Slots {
			if !dayExists || !daySlots[userSlot] {
				var ce exceptions.ConflictErrors
				ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
				return &ce
			}
		}
	}

	return nil
}

func (rs *RequestService) validateRequestPets(user *entities.User, petSitter *entities.PetSitter, petIDs []uint) ([]entities.Pet, error) {
	pets, err := rs.petService.GetPetsInUser(user.Pets, petIDs)
	if err != nil {
		return nil, err
	}

	if err := rs.petSitterService.ValidatePets(pets, petSitter.PetKinds); err != nil {
		return nil, err
	}

	return rs.makeRequestPets(pets)
}

func (rs *RequestService) makeRequestPets(userPets []entities.Pet) ([]entities.Pet, error) {
	pets := make([]entities.Pet, 0)

	for _, pet := range userPets {
		requestPet := &entities.Pet{
			UserID:     pet.UserID,
			Name:       pet.Name,
			Kind:       pet.Kind,
			Species:    pet.Species,
			BirthDate:  pet.BirthDate,
			IsAdult:    pet.IsAdult,
			Gender:     pet.Gender,
			Weight:     pet.Weight,
			PictureKey: pet.PictureKey,
			AboutPet:   pet.AboutPet,
			Type:       bootstrap.Run().Constants.EntityConstants.Request,
		}

		pets = append(pets, *requestPet)
	}

	return pets, nil
}

func (rs *RequestService) calculateTotalPrice(servicesEntity *entities.Service, slots []entities.CalendarSlot, petCount int) uint {
	var totalPrice uint = 0

	servicePrice := servicesEntity.Price * uint(petCount)
	totalPrice += servicePrice

	for _, slot := range slots {
		slotHours := uint(len(slot.Slots) / 2)
		totalPrice *= slotHours
	}

	return totalPrice
}

func (rs *RequestService) makeCalendarSlots(calendarSlots []request.RequestCalendarSlotRequest) []entities.CalendarSlot {
	slots := make([]entities.CalendarSlot, 0)

	for _, slot := range calendarSlots {
		calendarSlot := &entities.CalendarSlot{
			Date:   slot.Date,
			Slots:  slot.Slots,
			Status: enums.Booked,
		}
		slots = append(slots, *calendarSlot)
	}

	return slots
}

func (rs *RequestService) FindRequestByID(id uint) (*entities.Request, error) {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	foundRequest, err := requestRepo.GetRequestByID(id)
	if err != nil {
		return nil, err
	}
	if foundRequest == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	return foundRequest, nil
}

func (rs *RequestService) PreloadFields(request *entities.Request, fields []string) error {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	return requestRepo.PreloadFields(request, fields)
}
