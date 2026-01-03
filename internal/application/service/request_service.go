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
	domainpostgres "hona/backend/internal/domain/ports/postgres"
	"hona/backend/internal/infrastructure/communication/mail"
	"hona/backend/internal/infrastructure/persistence/repository/postgres"
	"log"
	"sort"
	"time"
)

type RequestService struct {
	userService      usecase.UserService
	addressService   usecase.AddressService
	petService       usecase.PetService
	petSitterService usecase.PetSitterService
	walletService    usecase.WalletService
	unitOfWork       ports.UnitOfWork
	emailService     *mail.EmailService
}

type RequestServiceDeps struct {
	UserService      usecase.UserService
	AddressService   usecase.AddressService
	PetService       usecase.PetService
	PetSitterService usecase.PetSitterService
	WalletService    usecase.WalletService
	UnitOfWork       ports.UnitOfWork
	EmailService     *mail.EmailService
}

func NewRequestService(deps RequestServiceDeps) *RequestService {
	return &RequestService{
		userService:      deps.UserService,
		unitOfWork:       deps.UnitOfWork,
		addressService:   deps.AddressService,
		petService:       deps.PetService,
		petSitterService: deps.PetSitterService,
		walletService:    deps.WalletService,
		emailService:     deps.EmailService,
	}
}

func (rs *RequestService) CreateRequest(info request.CreateRequestRequest) error {
	petSitter, err := rs.petSitterService.GetPetSitterByUserID(info.PetSitterUserID)
	if err != nil {
		return err
	}

	if petSitter.Status != enums.PSS_Active {
		return exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
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

	user.Pets, err = rs.petService.FindUserPetsByID(user.ID)
	if err != nil {
		return err
	}

	pets, err := rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	if info.AddressInfo != nil {
		madeAddress, err := rs.addressService.CreateAddressEntity(*info.AddressInfo)
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
		UserID:        info.UserID,
		PetSitterID:   petSitter.ID,
		Status:        enums.Pending,
		Chat:          entities.Chat{},
		TransferID:    nil,
		CalendarSlots: calendarSlots,
		Pets:          pets,
		TotalPrice:    totalPrice,
		Notes:         info.Notes,
		Comment:       nil,
		Address:       *address,
		Service:       *serviceEntity,
	}

	rs.sendNewRequestEmail(petSitter.UserID)

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

	pets, err := rs.petService.FindUserPetsByID(foundUser.ID)
	if err != nil {
		return nil, err
	}

	petsData, err := rs.petService.GetPetsBasicDataResponse(pets)
	if err != nil {
		return nil, err
	}

	if len(petsData) == 0 {
		return nil, exceptions.NewAccessDeniedError("first add a pet")
	}

	petSitter, err := rs.petSitterService.GetPetSitterByUserID(info.PetSitterUserID)
	if err != nil {
		return nil, err
	}

	filteredPetsData := rs.filterPetsByPetSitterKinds(petsData, petSitter)

	if petSitter.Status != enums.PSS_Active {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}

	petSitterUser, err := rs.userService.FindUserByID(petSitter.UserID)
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

	if len(freeSlots) == 0 {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}

	servicesData, err := rs.petSitterService.GetAvailableServicesResponse(petSitter)
	if err != nil {
		return nil, err
	}

	if len(servicesData) == 0 {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.PetSitter)
	}

	return &request.CreateRequestInfoResponse{
		Services:           servicesData,
		Addresses:          addresses,
		Pets:               filteredPetsData,
		FreeCalendarSlots:  freeSlots,
		PetSitterFirstName: petSitterUser.FirstName,
		PetSitterLastName:  petSitterUser.LastName,
	}, nil
}

func (rs *RequestService) filterPetsByPetSitterKinds(petsData []pet.PetBasicDataResponse, petSitter *entities.PetSitter) []pet.PetBasicDataResponse {
	filteredPetsData := make([]pet.PetBasicDataResponse, 0, len(petsData))
	for _, petData := range petsData {
		for _, kind := range petSitter.PetKinds {
			if petData.Kind == kind.String() {
				filteredPetsData = append(filteredPetsData, petData)
				break
			}
		}
	}

	return filteredPetsData
}

func (rs *RequestService) EditRequest(info request.EditRequestRequest) error {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return err
	}

	if foundRequest.Status != enums.Pending && foundRequest.Status != enums.Conflict {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterID)
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

	user.Pets, err = rs.petService.FindUserPetsByID(user.ID)
	if err != nil {
		return err
	}

	user.Pets, err = rs.petService.FindUserPetsByID(user.ID)
	if err != nil {
		return err
	}

	pets, err := rs.validateRequestPets(user, petSitter, info.PetIDs)
	if err != nil {
		return err
	}

	var address *entities.Address
	if info.AddressInfo != nil {
		madeAddress, err := rs.addressService.CreateAddressEntity(*info.AddressInfo)
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

	rs.sendEditRequestEmail(user, petSitter.UserID)

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	err = requestRepo.DeletePetsByRequestID(foundRequest.ID)
	if err != nil {
		return err
	}
	err = requestRepo.DeleteCalendarSlotsByRequestID(foundRequest.ID)
	if err != nil {
		return err
	}
	err = requestRepo.EditRequest(foundRequest)
	if err != nil {
		return err
	}

	return nil
}

func (rs *RequestService) CancelRequest(info request.CancelRequestRequest) error {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return err
	}

	err = rs.PreloadFields(foundRequest, []string{"CalendarSlots"})
	if err != nil {
		return err
	}
	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterID)
	if err != nil {
		return err
	}

	if foundRequest.UserID != info.UserID && petSitter.UserID != info.UserID {
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

	err = rs.petSitterService.PreloadFields(petSitter, []string{"Schedule"})
	if err != nil {
		return err
	}

	foundRequest.Status = enums.Canceled

	err = rs.petSitterService.AutoUpdateSlots(petSitter, foundRequest.CalendarSlots, false)
	if err != nil {
		return err
	}
	if info.UserID == petSitter.UserID {
		rs.SendPetOwnerRequestCancelEmail(foundRequest.UserID, petSitter.UserID)
	} else {
		rs.SendPetSitterRequestCancelEmail(foundRequest.UserID, petSitter.UserID)
	}
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	return requestRepo.EditRequest(foundRequest)
}

func (rs *RequestService) PayRequest(info request.PayRequestRequest) error {
	return rs.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		requestRepo := rf.RequestRepository()
		transferRepo := rf.TransferRepository()

		foundRequest, err := rs.loadPayableRequest(requestRepo, info)
		if err != nil {
			return err
		}

		petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterID)
		if err != nil {
			return err
		}

		senderWallet, receiverWallet, err := rs.walletService.TransferInTransaction(rf, foundRequest.UserID, petSitter.UserID, foundRequest.TotalPrice)
		if err != nil {
			return err
		}

		transfer, err := rs.createTransfer(transferRepo, senderWallet.ID, receiverWallet.ID, foundRequest.TotalPrice)
		if err != nil {
			return err
		}

		return rs.markRequestPaid(requestRepo, foundRequest, transfer.ID)
	})
}

func (rs *RequestService) loadPayableRequest(requestRepo domainpostgres.RequestRepository, info request.PayRequestRequest) (*entities.Request, error) {
	foundRequest, err := requestRepo.GetRequestByID(info.RequestID)
	if err != nil {
		return nil, err
	}
	if foundRequest == nil {
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Request)
	}
	if foundRequest.UserID != info.UserID {
		return nil, exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
	}
	if foundRequest.Status != enums.Accepted {
		return nil, exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
	}

	return foundRequest, nil
}

func (rs *RequestService) createTransfer(transferRepo domainpostgres.TransferRepository, senderWalletID, receiverWalletID, amount uint) (*entities.Transfer, error) {
	transfer := &entities.Transfer{
		ReceiverWalletID: receiverWalletID,
		SenderWalletID:   senderWalletID,
		Amount:           amount,
	}
	if err := transferRepo.CreateTransfer(transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (rs *RequestService) markRequestPaid(requestRepo domainpostgres.RequestRepository, foundRequest *entities.Request, transferID uint) error {
	foundRequest.Status = enums.Paid
	foundRequest.TransferID = &transferID
	return requestRepo.EditRequest(foundRequest)
}

func buildRequestStatusResponse(status enums.RequestStatus) request.RequestStatusResponse {
	return request.RequestStatusResponse{
		Num:  status,
		Name: status.String(),
	}
}

func (rs *RequestService) SearchRequests(info request.SearchRequestsRequest) ([]request.RequestListItemResponse, int64, error) {
	options := postgres.NewQueryOptions().WithPagination(info.Limit, info.Offset)
	if len(info.Filters) > 0 {
		options.WithFilters(info.Filters)
	}
	if len(info.Sorts) > 0 {
		options.WithSorting(info.Sorts)
	}

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	requests, total, err := requestRepo.SearchRequests(info.UserID, options)
	if err != nil {
		return nil, 0, err
	}

	requestUser, err := rs.userService.FindUserByID(info.UserID)
	if err != nil {
		return nil, 0, err
	}
	userPictureLink, err := rs.userService.GetUserPictureLink(requestUser)
	if err != nil {
		return nil, 0, err
	}

	for i := range requests {
		if err := rs.PreloadFields(&requests[i], []string{"Service", "CalendarSlots"}); err != nil {
			return nil, 0, err
		}
	}

	res := make([]request.RequestListItemResponse, len(requests))
	for i := range requests {
		req := requests[i]
		petSitter, err := rs.petSitterService.GetPetSitterByID(req.PetSitterID)
		if err != nil {
			return nil, 0, err
		}
		petSitterUser, err := rs.userService.FindUserByID(petSitter.UserID)
		if err != nil {
			return nil, 0, err
		}
		petSitterPictureLink, err := rs.userService.GetUserPictureLink(petSitterUser)
		if err != nil {
			return nil, 0, err
		}
		address, err := rs.addressService.FindRequestAddressByID(req.ID)
		if err != nil {
			return nil, 0, err
		}
		res[i] = request.RequestListItemResponse{
			RequestID:            req.ID,
			PetSitterUserID:      petSitter.UserID,
			PetSitterFirstName:   petSitterUser.FirstName,
			PetSitterLastName:    petSitterUser.LastName,
			PetSitterPictureLink: petSitterPictureLink,
			UserPictureLink:      userPictureLink,
			Service:              rs.petSitterService.GetServiceResponse(&req.Service),
			CalendarSlots:        rs.petSitterService.GetCalendarSlotsResponse(req.CalendarSlots),
			Address:              rs.addressService.GetUserAddressInfo(address),
			TotalPrice:           req.TotalPrice,
			Status:               buildRequestStatusResponse(req.Status),
			UpdatedAt:            req.UpdatedAt,
		}
	}

	return res, total, nil
}

func (rs *RequestService) SearchPetSitterRequests(info request.SearchPetSitterRequestsRequest) ([]request.RequestListItemResponse, int64, error) {
	petSitter, err := rs.petSitterService.GetPetSitterByUserID(info.PetSitterUserID)
	if err != nil {
		return nil, 0, err
	}
	petSitterUser, err := rs.userService.FindUserByID(petSitter.UserID)
	if err != nil {
		return nil, 0, err
	}
	petSitterPictureLink, err := rs.userService.GetUserPictureLink(petSitterUser)
	if err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().WithPagination(info.Limit, info.Offset)
	if len(info.Filters) > 0 {
		options.WithFilters(info.Filters)
	}
	if len(info.Sorts) > 0 {
		options.WithSorting(info.Sorts)
	}

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	requests, total, err := requestRepo.SearchRequestsByPetSitterID(petSitter.ID, options)
	if err != nil {
		return nil, 0, err
	}

	for i := range requests {
		if err := rs.PreloadFields(&requests[i], []string{"Service", "CalendarSlots"}); err != nil {
			return nil, 0, err
		}
	}

	res := make([]request.RequestListItemResponse, len(requests))
	for i := range requests {
		req := requests[i]
		address, err := rs.addressService.FindRequestAddressByID(req.ID)
		if err != nil {
			return nil, 0, err
		}
		requestUser, err := rs.userService.FindUserByID(req.UserID)
		if err != nil {
			return nil, 0, err
		}
		userPictureLink, err := rs.userService.GetUserPictureLink(requestUser)
		if err != nil {
			return nil, 0, err
		}
		res[i] = request.RequestListItemResponse{
			RequestID:            req.ID,
			PetSitterUserID:      petSitter.UserID,
			PetSitterFirstName:   petSitterUser.FirstName,
			PetSitterLastName:    petSitterUser.LastName,
			PetSitterPictureLink: petSitterPictureLink,
			UserPictureLink:      userPictureLink,
			Service:              rs.petSitterService.GetServiceResponse(&req.Service),
			CalendarSlots:        rs.petSitterService.GetCalendarSlotsResponse(req.CalendarSlots),
			Address:              rs.addressService.GetUserAddressInfo(address),
			TotalPrice:           req.TotalPrice,
			Status:               buildRequestStatusResponse(req.Status),
			UpdatedAt:            req.UpdatedAt,
		}
	}

	return res, total, nil
}

func (rs *RequestService) GetRequestFullData(info request.GetRequestFullDataRequest) (*request.RequestFullDataResponse, error) {
	foundRequest, err := rs.FindRequestByID(info.RequestID)
	if err != nil {
		return nil, err
	}

	err = rs.PreloadFields(foundRequest, []string{"CalendarSlots", "Service", "Comment"})
	if err != nil {
		return nil, err
	}

	pets, err := rs.petService.FindRequestPetsByID(foundRequest.ID)
	if err != nil {
		return nil, err
	}

	address, err := rs.addressService.FindRequestAddressByID(foundRequest.ID)
	if err != nil {
		return nil, err
	}

	petsData, err := rs.petService.GetPetsBasicDataResponse(pets)
	if err != nil {
		return nil, err
	}

	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterID)
	if err != nil {
		return nil, err
	}

	petSitterUser, err := rs.userService.FindUserByID(petSitter.UserID)
	if err != nil {
		return nil, err
	}

	requestUser, err := rs.userService.FindUserByID(foundRequest.UserID)
	if err != nil {
		return nil, err
	}
	userPictureLink, err := rs.userService.GetUserPictureLink(requestUser)
	if err != nil {
		return nil, err
	}

	commentResponse, err := buildCommentResponse(rs.userService, requestUser, foundRequest.Comment)
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
		RequestID:          foundRequest.ID,
		PetSitterUserID:    petSitter.UserID,
		PetSitterFirstName: petSitterUser.FirstName,
		PetSitterLastName:  petSitterUser.LastName,
		UserFirstName:      requestUser.FirstName,
		UserLastName:       requestUser.LastName,
		UserPictureLink:    userPictureLink,
		Service:            rs.petSitterService.GetServiceResponse(&foundRequest.Service),
		Pets:               petsData,
		Address:            rs.addressService.GetUserAddressInfo(address),
		Notes:              foundRequest.Notes,
		TotalPrice:         foundRequest.TotalPrice,
		Comment:            commentResponse,
		Status:             buildRequestStatusResponse(foundRequest.Status),
		TransferID:         foundRequest.TransferID,
		CalendarSlots:      rs.petSitterService.GetCalendarSlotsResponse(foundRequest.CalendarSlots),
		UpdatedAt:          foundRequest.UpdatedAt,
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
	petSitter, err := rs.petSitterService.GetPetSitterByID(foundRequest.PetSitterID)
	if err != nil {
		return err
	}
	if foundRequest.Status != enums.Pending {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}
	if petSitter.UserID != info.UserID {
		err = exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
		return err
	}

	if foundRequest.UpdatedAt.After(info.GetTime) {
		var ce exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Request, bootstrap.Run().Constants.ErrorTags.OldInfo)
		return &ce
	}

	if info.Accept {
		foundRequest.Status = enums.Accepted

		err = rs.validateCalendarSlots(petSitter.Schedule, foundRequest.CalendarSlots)
		if err != nil {
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.CalendarSlot, bootstrap.Run().Constants.ErrorTags.CalendarConflict)
		}

		err = rs.petSitterService.AutoUpdateSlots(petSitter, foundRequest.CalendarSlots, true)
		if err != nil {
			return err
		}

		rs.sendAcceptRequestEmail(foundRequest.UserID)
	} else {
		foundRequest.Status = enums.Dismissed

		rs.sendDeclineRequestEmail(foundRequest.UserID)
	}

	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	return requestRepo.EditRequest(foundRequest)
}

func (rs *RequestService) validateCalendarSlots(petSitterSlots []entities.CalendarSlot, requestSlots []entities.CalendarSlot) error {
	availableSlots := rs.petSitterService.GetFreeMap(petSitterSlots)

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

	if err := rs.petSitterService.ValidatePets(pets, []enums.PetKind(petSitter.PetKinds)); err != nil {
		return nil, err
	}

	return rs.makeRequestPets(pets)
}

func (rs *RequestService) makeRequestPets(userPets []entities.Pet) ([]entities.Pet, error) {
	pets := make([]entities.Pet, len(userPets))

	for i, pet := range userPets {
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

		pets[i] = *requestPet
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

	return uint(totalPrice + totalPrice/10)
}

func (rs *RequestService) makeCalendarSlots(calendarSlots []request.RequestCalendarSlotRequest) []entities.CalendarSlot {
	slots := make([]entities.CalendarSlot, len(calendarSlots))

	for i, slot := range calendarSlots {
		calendarSlot := &entities.CalendarSlot{
			Date:   slot.Date,
			Slots:  entities.Slots(slot.Slots),
			Status: enums.Booked,
		}
		slots[i] = *calendarSlot
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
		return nil, exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Request)
	}

	return foundRequest, nil
}

func (rs *RequestService) PreloadFields(request *entities.Request, fields []string) error {
	requestRepo := rs.unitOfWork.Factory().RequestRepository()
	return requestRepo.PreloadFields(request, fields)
}

func (rs *RequestService) sendNewRequestEmail(id uint) {
	petSitterUser, err := rs.userService.FindUserByID(id)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		Year int
	}{
		Year: time.Now().Year(),
	}
	err = rs.emailService.SendEmail(petSitterUser.Email, "New Request Received", bootstrap.Run().Constants.TemplatesPath.NewRequest, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) sendEditRequestEmail(user *entities.User, id uint) {
	petSitterUser, err := rs.userService.FindUserByID(id)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		RequesterName string
		Year          int
	}{
		RequesterName: user.FirstName,
		Year:          time.Now().Year(),
	}
	err = rs.emailService.SendEmail(petSitterUser.Email, "Request Edited", bootstrap.Run().Constants.TemplatesPath.RequestEdited, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) SendPetOwnerRequestCancelEmail(userID, petSitterUserID uint) {
	user, err := rs.userService.FindUserByID(userID)
	if err != nil {
		log.Println(err)
	}
	petSitterUser, err := rs.userService.FindUserByID(petSitterUserID)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		RequesterName string
		SitterName    string
		Year          int
	}{
		RequesterName: user.FirstName,
		SitterName:    petSitterUser.FirstName + " " + petSitterUser.LastName,
		Year:          time.Now().Year(),
	}
	err = rs.emailService.SendEmail(user.Email, "Request Canceled", bootstrap.Run().Constants.TemplatesPath.PetOwnerRequestCancel, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) SendPetSitterRequestCancelEmail(userID, petSitterUserID uint) {
	user, err := rs.userService.FindUserByID(userID)
	if err != nil {
		log.Println(err)
	}
	petSitterUser, err := rs.userService.FindUserByID(petSitterUserID)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		RequesterName string
		Year          int
	}{
		RequesterName: user.FirstName,
		Year:          time.Now().Year(),
	}
	err = rs.emailService.SendEmail(petSitterUser.Email, "Request Canceled", bootstrap.Run().Constants.TemplatesPath.PetSitterRequestCancel, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) sendAcceptRequestEmail(id uint) {
	user, err := rs.userService.FindUserByID(id)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		RequesterName string
		Year          int
	}{
		RequesterName: user.FirstName,
		Year:          time.Now().Year(),
	}
	err = rs.emailService.SendEmail(user.Email, "Request Accepted", bootstrap.Run().Constants.TemplatesPath.RequestAccepted, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) sendDeclineRequestEmail(id uint) {
	user, err := rs.userService.FindUserByID(id)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		RequesterName string
		Year          int
	}{
		RequesterName: user.FirstName,
		Year:          time.Now().Year(),
	}
	err = rs.emailService.SendEmail(user.Email, "Request Declined", bootstrap.Run().Constants.TemplatesPath.RequestDeclined, data)
	if err != nil {
		log.Println(err)
	}
}

func (rs *RequestService) EnsureRequestIsFinished(request *entities.Request) error {
	if request.Status != enums.Finished {
		return exceptions.NewAccessDeniedError(bootstrap.Run().Constants.ErrorTags.ForbiddenStatus)
	}
	return nil
}
