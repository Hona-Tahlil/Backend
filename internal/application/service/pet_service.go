package service

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/usecase"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	domainstorage "hona/backend/internal/domain/storage"
	"log"
	"time"
)

type PetService struct {
	unitOfWork  ports.UnitOfWork
	storage     domainstorage.Storage
	userService usecase.UserService
}

func NewPetService(unitOfWork ports.UnitOfWork, storage domainstorage.Storage, userService usecase.UserService) *PetService {
	return &PetService{
		unitOfWork:  unitOfWork,
		storage:     storage,
		userService: userService,
	}
}

func (ps *PetService) AddPet(info pet.AddPetRequest) error {
	err := ps.validateSpecies(info.Species, info.Kind)
	if err != nil {
		return err
	}
	isAdult, err := ps.validateBirthDate(info.BirthDate, info.IsAdult)
	if err != nil {
		return err
	}

	info.IsAdult = isAdult
	var profileKey *string
	if info.ProfilePic != nil {
		profileKeyValue := ps.getStorageKey(info.Name, info.UserID)
		profileKey = &profileKeyValue
		if err := ps.storage.UploadFile(enums.PetProfilePic, *profileKey, info.ProfilePic); err != nil {
			return err
		}
	}

	_, err = ps.findPet(info.Name, info.UserID)
	if err == nil {
		var ce exceptions.ValidationErrors
		ce.AddError(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.DuplicateName)
		return &ce
	} else if _, ok := err.(*exceptions.NotFoundError); !ok {
		return err
	}

	pet := &entities.Pet{
		UserID:     info.UserID,
		Name:       info.Name,
		Kind:       info.Kind,
		Species:    info.Species,
		BirthDate:  info.BirthDate,
		IsAdult:    info.IsAdult,
		Gender:     info.Gender,
		Weight:     info.Weight,
		PictureKey: profileKey,
		AboutPet:   info.AboutPet,
	}
	petRepo := ps.unitOfWork.Factory().PetRepository()
	err = petRepo.CreatePet(pet)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PetService) UpdatePet(info pet.UpdatePetRequest) error {
	foundPet, err := ps.findPetByID(info.ID)
	if err != nil {
		return err
	}

	err = ps.validateSpecies(info.Species, info.Kind)
	if err != nil {
		return err
	}

	isAdult, err := ps.validateBirthDate(info.BirthDate, info.IsAdult)
	if err != nil {
		return err
	}

	info.IsAdult = isAdult
	var profileKey *string
	if info.ProfilePic != nil {
		profileKeyValue := ps.getStorageKey(info.Name, foundPet.UserID)
		profileKey = &profileKeyValue
		if err := ps.storage.UploadFile(enums.PetProfilePic, *profileKey, info.ProfilePic); err != nil {
			return err
		}
	}

	_, err = ps.findPet(info.Name, foundPet.UserID)
	if err == nil {
		var ce exceptions.ConflictErrors
		ce.Add(bootstrap.Run().Constants.ErrorFields.Pet, bootstrap.Run().Constants.ErrorTags.DuplicateName)
		return &ce
	} else if _, ok := err.(*exceptions.NotFoundError); !ok {
		return err
	}

	foundPet.Name = info.Name
	foundPet.AboutPet = info.AboutPet
	foundPet.BirthDate = info.BirthDate
	foundPet.Gender = info.Gender
	foundPet.IsAdult = info.IsAdult
	foundPet.Kind = info.Kind
	foundPet.PictureKey = profileKey
	foundPet.Species = info.Species
	foundPet.Weight = info.Weight

	petRepo := ps.unitOfWork.Factory().PetRepository()
	err = petRepo.UpdatePet(foundPet)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PetService) RemovePet(info pet.RemovePetRequest) error {
	foundPet, err := ps.findPetByID(info.ID)
	if err != nil {
		return err
	}
	profileKeyValue := ps.getStorageKey(foundPet.Name, foundPet.UserID)
	if err = ps.storage.DeleteObject(enums.PetProfilePic, profileKeyValue); err != nil {
		log.Println(err)
	}
	petRepo := ps.unitOfWork.Factory().PetRepository()
	err = petRepo.RemovePet(foundPet)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PetService) GetPetsBasicData(info pet.GetPetsBasicDataRequest) ([]pet.PetBasicDataResponse, error) {
	r := make([]pet.PetBasicDataResponse, 0)
	user, err := ps.userService.FindUserByID(info.UserID)
	if err != nil {
		return nil, err
	}
	petRepo := ps.unitOfWork.Factory().PetRepository()
	err = petRepo.PreloadUserPets(user)
	if err != nil {
		return nil, err
	}
	for _, pet := range user.Pets {
		res, err := ps.getPetBasicDataResponse(&pet)
		if err != nil {
			return nil, err
		}
		r = append(r, *res)
	}
	return r, nil
}

func (ps *PetService) GetPetFullData(info pet.GetPetFullDataRequest) (*pet.PetFullDataResponse, error) {
	foundPet, err := ps.findPetByID(info.ID)
	if err != nil {
		return nil, err
	}
	profileKeyValue := ps.getStorageKey(foundPet.Name, foundPet.UserID)
	link, err := ps.storage.GetPresignedURL(enums.PetProfilePic, profileKeyValue, time.Minute*15)
	if err != nil {
		log.Println(err)
	}

	return &pet.PetFullDataResponse{
		ID:          foundPet.ID,
		Name:        foundPet.Name,
		Kind:        foundPet.Kind,
		Species:     foundPet.Species,
		Gender:      foundPet.Gender,
		PictureLink: link,
		BirthDate:   foundPet.BirthDate,
		IsAdult:     foundPet.IsAdult,
		Weight:      foundPet.Weight,
		AboutPet:    foundPet.AboutPet,
	}, nil
}

func (ps *PetService) getPetBasicDataResponse(petEntity *entities.Pet) (*pet.PetBasicDataResponse, error) {
	profileKeyValue := ps.getStorageKey(petEntity.Name, petEntity.UserID)
	link, err := ps.storage.GetPresignedURL(enums.PetProfilePic, profileKeyValue, time.Minute*15)
	if err != nil {
		log.Println(err)
	}
	return &pet.PetBasicDataResponse{
		ID:          petEntity.ID,
		Name:        petEntity.Name,
		Kind:        petEntity.Kind,
		Species:     petEntity.Species,
		Gender:      petEntity.Gender,
		PictureLink: link,
	}, nil

}

func (ps *PetService) findPet(name string, userID uint) (*entities.Pet, error) {
	petRepo := ps.unitOfWork.Factory().PetRepository()
	foundPet, err := petRepo.FindPet(name, userID)
	if err != nil {
		return nil, err
	}
	if foundPet == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Pet)
		return nil, NotFoundError
	}
	return foundPet, nil
}

func (ps *PetService) findPetByID(id uint) (*entities.Pet, error) {
	petRepo := ps.unitOfWork.Factory().PetRepository()
	foundPet, err := petRepo.FindPetByID(id)
	if err != nil {
		return nil, err
	}
	if foundPet == nil {
		NotFoundError := exceptions.NewNotFoundError(bootstrap.Run().Constants.ErrorFields.Pet)
		return nil, NotFoundError
	}
	return foundPet, nil
}

func (ps *PetService) validateBirthDate(birthDate *time.Time, isAdultInput bool) (isAdult bool, err error) {
	if birthDate != nil {
		age := time.Now().Year() - birthDate.Year()
		if time.Now().YearDay() < birthDate.YearDay() {
			age--
		}
		if age > 2 {
			isAdult = true
		} else {
			isAdult = false
		}
		if age > 100 || time.Now().Before(*birthDate) {
			var ce exceptions.ValidationErrors
			ce.AddError(bootstrap.Run().Constants.ErrorFields.BirthDate, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			err = &ce
		}

		return
	}
	isAdult = isAdultInput
	return
}

func (ps *PetService) getStorageKey(name string, userID uint) string {
	return "profile-" + name + "-" + fmt.Sprint(userID)
}

func (ps *PetService) validateSpecies(species enums.Species, kind enums.PetKind) error {
	var invalid = func() error {
		var ce exceptions.ValidationErrors
		ce.AddError(
			bootstrap.Run().Constants.ErrorFields.Species,
			bootstrap.Run().Constants.ErrorTags.UnacceptableInput,
		)
		return &ce
	}

	if species == enums.Other {
		return nil
	}
	switch kind {
	case enums.Dog:
		if species < 1 || species > 14 {
			return invalid()
		}
	case enums.Cat:
		if species < 15 || species > 22 {
			return invalid()
		}
	case enums.Bird:
		if species < 23 || species > 30 {
			return invalid()
		}
	case enums.Fish:
		if species < 31 || species > 36 {
			return invalid()
		}
	case enums.Rodent:
		if species < 37 || species > 43 {
			return invalid()
		}
	case enums.Rabbit:
		if species < 44 || species > 47 {
			return invalid()
		}
	case enums.Reptile:
		if species < 48 || species > 55 {
			return invalid()
		}
	case enums.Amphibian:
		if species < 56 || species > 58 {
			return invalid()
		}
	case enums.Ferret:
		if species != 59 {
			return invalid()
		}
	case enums.Horse:
		if species != 60 {
			return invalid()
		}
	case enums.Hedgehog:
		if species != 61 {
			return invalid()
		}
	case enums.MiniPig:
		if species != 62 {
			return invalid()
		}
	case enums.Insect:
		if species < 63 || species > 64 {
			return invalid()
		}
	case enums.Arachnid:
		if species != 65 {
			return invalid()
		}
	case enums.HermitCrab:
		if species != 66 {
			return invalid()
		}
	}

	return nil
}

func (ps *PetService) GetAllPetKinds() []pet.PetKindResponse {
	kinds := enums.GetAllPetKinds()
	res := make([]pet.PetKindResponse, 0)
	for _, kind := range kinds {
		res = append(res, pet.PetKindResponse{
			Num:  kind,
			Name: kind.String(),
		})
	}
	return res
}

func (ps *PetService) GetPetKindSpecies(info pet.GetPetKindSpecies) []pet.PetSpeciesResponse {
	species := enums.GetSpeciesByKind(info.Num)
	species = append(species, enums.Other)
	res := make([]pet.PetSpeciesResponse, 0)
	for _, s := range species {
		res = append(res, pet.PetSpeciesResponse{
			Num:  s,
			Name: s.String(),
		})
	}
	return res
}
