package service

import (
	"fmt"
	"hona/backend/bootstrap"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"hona/backend/internal/domain/exceptions"
	"hona/backend/internal/domain/ports"
	"hona/backend/internal/infrastructure/storage"
	"time"
)

type PetService struct {
	unitOfWork ports.UnitOfWork
	storage    storage.S3Storage
}

func NewPetService(unitOfWork ports.UnitOfWork, storage storage.S3Storage) *PetService {
	return &PetService{
		unitOfWork: unitOfWork,
		storage:    storage,
	}
}

func (ps *PetService) AddPet(info pet.AddPetRequest) error {
	// TODO: func to validate correct species
	isAdult, err := ps.validateBirthDate(info.BirthDate, info.IsAdult, info.Kind)
	if err != nil {
		return err
	}

	info.IsAdult = isAdult
	var profileKey *string
	if info.ProfilePic != nil {
		profileKeyValue := "profile-" + info.Name + "-" + fmt.Sprint(info.UserID)
		profileKey = &profileKeyValue
		if err := ps.storage.UploadFile(enums.PetProfilePic, *profileKey, info.ProfilePic); err != nil {
			return err
		}
	}

	// TODO: validate duplicate name

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
	// TODO: func to validate correct species
	isAdult, err := ps.validateBirthDate(info.BirthDate, info.IsAdult, info.Kind)
	if err != nil {
		return err
	}

	info.IsAdult = isAdult
	var profileKey *string
	if info.ProfilePic != nil {
		profileKeyValue := "profile-" + info.Name + "-" + fmt.Sprint(info.UserID)
		profileKey = &profileKeyValue
		if err := ps.storage.UploadFile(enums.PetProfilePic, *profileKey, info.ProfilePic); err != nil {
			return err
		}
	}

	// TODO: validate duplicate name

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

func (ps *PetService) findPet(name string, userID uint) (*entities.Pet, error) {
	// TODO: implement
	return nil, nil
}

func (ps *PetService) validateBirthDate(birthDate *time.Time, isAdultInput bool, kind enums.PetKind) (isAdult bool, err error) {
	if birthDate != nil {
		// TODO: further implemention based on kind
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
			var ce exceptions.ConflictErrors
			ce.Add(bootstrap.Run().Constants.ErrorFields.BirthDate, bootstrap.Run().Constants.ErrorTags.UnacceptableInput)
			err = ce
		}

		return
	}
	isAdult = isAdultInput
	return
}
