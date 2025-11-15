package service

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/pet"
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
	cityService         usecase.CityService
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

func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	petSitterUser, err := rs.userService.FindUserByID(info.PetSitterUserID)
	if err != nil {
		return err
	}
	// TODO: move this to user service
	userRepo := rs.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(petSitterUser)
	if err != nil {
		return err
	}
	if petSitterUser.PetSitter == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
	}
	petSitter := petSitterUser.PetSitter
	if !petSitter.IsVerified {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
	}

	petSitterRepo := rs.unitOfWork.Factory().PetSitterRepository()
	err = petSitterRepo.PreloadSchedule(petSitter)
	if err != nil {
		return err
	}

	err = rs.validateCalendarSlots(petSitter.Schedule, info.CalenderSlots)
	if err != nil {
		return err
	}

	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	if !user.IsEmailVerified {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.User, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ve
	}

	err = rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	pets, err := rs.makeRequestPets(info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	// TODO: logic in address service
	if info.AddressInfo != nil {
		province, err := rs.provinceService.FindProvinceByName(info.AddressInfo.ProvinceName)
		if err != nil {
			return err
		}

		city, err := rs.cityService.FindCityByNameInProvince(info.AddressInfo.CityName, province)
		if err != nil {
			return err
		}

		address = &entities.Address{
			Province:      *province,
			City:          *city,
			StreetAddress: info.AddressInfo.StreetAddress,
			HouseNumber:   info.AddressInfo.HouseNumber,
			Unit:          info.AddressInfo.Unit,
			PostalCode:    info.AddressInfo.PostalCode,
		}
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

	err = rs.validateRequestService(petSitter, info.ServiceID)
	if err != nil {
		return err
	}

	serviceEntity, err := rs.makeRequestService(info.ServiceID)
	if err != nil {
		return err
	}

	totalPrice := rs.calculateTotalPrice(serviceEntity, petSitter.Schedule, len(pets))

	calendarSlots := rs.makeCalendarSlots(info.CalenderSlots)

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

	petsData, err := rs.petService.GetPetsBasicData(pet.GetPetsBasicDataRequest{
		UserID: info.UserID,
	})
	if err != nil {
		return nil, err
	}

	freeSlots, err := rs.petSitterService.GetPetSitterFreeSlotsResponse(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	servicesData, err := rs.petSitterService.GetServicesResponse(info.PetSitterUserID)
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
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	foundRequest, err := requestRepo.GetRequestByID(info.RequestID)
	if err != nil {
		return err
	}
	if foundRequest == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
	}
	if foundRequest.Status != enums.Pending {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	petSitterUser, err := rs.userService.FindUserByID(foundRequest.PetSitterUserID)
	if err != nil {
		return err
	}
	// TODO: same as createRequest
	userRepo := rs.unitOfWork.Factory().UserRepository()
	err = userRepo.PreloadPetSitter(petSitterUser)
	if err != nil {
		return err
	}
	if petSitterUser.PetSitter == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
	}
	petSitter := petSitterUser.PetSitter
	if !petSitter.IsVerified {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.PetSitter, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
	}

	// TODO: move logic to the user service
	user, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return err
	}
	if !user.IsEmailVerified {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.User, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ve
	}

	err = rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	pets, err := rs.makeRequestPets(info.PetIDs)
	if err != nil {
		return err
	}

	// TODO: same as createRequest
	var address *entities.Address
	if info.AddressInfo != nil {
		province, err := rs.provinceService.FindProvinceByName(info.AddressInfo.ProvinceName)
		if err != nil {
			return err
		}

		city, err := rs.cityService.FindCityByNameInProvince(info.AddressInfo.CityName, province)
		if err != nil {
			return err
		}

		address = &entities.Address{
			Province:      *province,
			City:          *city,
			StreetAddress: info.AddressInfo.StreetAddress,
			HouseNumber:   info.AddressInfo.HouseNumber,
			Unit:          info.AddressInfo.Unit,
			PostalCode:    info.AddressInfo.PostalCode,
		}
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

	err = rs.validateRequestService(petSitter, info.ServiceID)
	if err != nil {
		return err
	}

	serviceEntity, err := rs.makeRequestService(info.ServiceID)
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

	requestRepo.EditRequest(foundRequest)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RequestService) CancelRequest(info request.CancelRequestRequest) error {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	// TODO: move this to dedicated method
	foundRequest, err := requestRepo.GetRequestByID(info.RequestID)
	if err != nil {
		return nil
	}
	if foundRequest == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
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

	requestRepo.EditRequest(foundRequest)

	return nil
}

func (rs *RequestService) GetRequestFullData(info request.GetRequestFullDataRequest) (*request.RequestFullDataResponse, error) {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	foundRequest, err := requestRepo.GetRequestByID(info.RequestID)
	if err != nil {
		return nil, err
	}
	if foundRequest == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.NotFound)
		return nil, &ve
	}

	petsData, err := rs.petService.GetPetsBasicDataResponse(foundRequest.Pets)
	if err != nil {
		return nil, err
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

// TODO: handle conflict here and in gets
func (rs *RequestService) RespondToRequest(info request.RespondToRequestRequest) error {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	foundRequest, err := requestRepo.GetRequestByID(info.RequestID)
	if err != nil {
		return err
	}
	if foundRequest == nil {
		var ve exceptions.ValidationErrors
		ve.AddError(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.NotFound)
		return &ve
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
	} else {
		foundRequest.Status = enums.Dismissed
	}

	requestRepo.EditRequest(foundRequest)

	return nil
}

func (rs *RequestService) validateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []request.RequestCalendarSlotRequest) error {
	for _, req := range requestSlots {
		for _, userSlot := range req.Slots {
			flag := false
			// TODO: map
			for _, petSitterSlot := range petSitterSlots {
				if petSitterSlot.Date.Equal(req.Date) && petSitterSlot.Status == enums.Free {
					for _, sitterSlot := range petSitterSlot.Slots {
						if sitterSlot == userSlot {
							flag = true
							break
						}
					}
				}
			}
			if !flag {
				var ce exceptions.ConflictErrors
				ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
				return &ce
			}
		}
	}

	return nil
}

func (rs *RequestService) validateRequestPets(user *entities.User, petSitter *entities.PetSitter, petIDs []uint) error {
	pets := make([]entities.Pet, 0)
	userRepo := rs.unitOfWork.Factory().UserRepository()
	err := userRepo.PreloadPets(user)
	if err != nil {
		return err
	}
	// TODO: pet service
	for _, petID := range petIDs {
		flag := false
		for _, pet := range user.Pets {
			if pet.ID == petID {
				flag = true
				pets = append(pets, pet)
				break
			}
		}
		if !flag {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			return &ce
		}
	}

	for _, pet := range pets {
		flag := false
		for _, kind := range petSitter.PetKinds {
			if kind == pet.Kind {
				flag = true
				break
			}
		}
		if !flag {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			return &ce
		}
	}

	return nil
}

func (rs *RequestService) makeRequestPets(petIDs []uint) ([]entities.Pet, error) {
	// TODO: use len and index instead of append
	pets := make([]entities.Pet, 0)

	for _, petID := range petIDs {
		pet, err := rs.petService.FindPetByID(petID)
		if err != nil {
			return nil, err
		}
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
			// TODO: constant
			Type: "request",
		}

		pets = append(pets, *requestPet)
	}

	return pets, nil
}

func (rs *RequestService) validateRequestService(petSitter *entities.PetSitter, serviceID uint) error {
	// TODO: Move this to pet sitter service
	flag := false
	for _, service := range petSitter.Services {
		if service.ID == serviceID {
			flag = true
			break
		}
	}
	if !flag {
		var ce exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Service, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
		return &ce
	}

	return nil
}

func (rs *RequestService) makeRequestService(serviceID uint) (*entities.Service, error) {
	service, err := rs.serviceService.FindServiceByID(serviceID)
	if err != nil {
		return nil, err
	}
	requestService := &entities.Service{
		PetSitterUserID: service.PetSitterUserID,
		Description:     service.Description,
		Price:           service.Price,
		Type:            service.Type,
		PetKinds:        service.PetKinds,
		// TODO: constant
		Kind: "request",
	}

	return requestService, nil
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
			Date:  slot.Date,
			Slots: slot.Slots,
			// TODO: status
		}
		slots = append(slots, *calendarSlot)
	}

	return slots
}
