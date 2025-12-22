package petsitter

import (
	"hona/backend/internal/application/dto/general"
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"time"
)

type AddressInfo struct {
	Province    enums.Province
	City        enums.City
	Address     string
	HouseNumber uint
	Unit        uint
	PostalCode  string
}

type GetPetSitterRequest struct {
	UserID uint
}

type SubmitPersonalInfoRequest struct {
	UserID      uint
	FirstName   string
	LastName    string
	Email       string
	Gender      enums.Gender
	BirthDate   *time.Time
	Phone       string
	Province    enums.Province
	City        enums.City
	Address     string
	HouseNumber uint
	Unit        uint
	PostalCode  string
}

type UploadDocumentsRequest struct {
	UserID           uint
	CertificateFiles []*multipart.FileHeader
	Files            []*multipart.FileHeader
}

type SubmitSkillsRequest struct {
	UserID   uint
	Bio      string
	PetKinds []enums.PetKind
	Services []enums.ServiceType
}

type GetPetSitterDetailsRequest struct {
	PetSitterUserID uint
}

type ChangePetSitterStatusRequest struct {
	PetSitterUserID uint
	Status          enums.PetSitterStatus
}

type SearchPetSittersRequest struct {
	Offset  int
	Limit   int
	Filters []general.Filter
	Sorts   []general.Sort
}

type AdminSearchPetSittersRequest struct {
	Offset  int
	Limit   int
	Filters []general.Filter
	Sorts   []general.Sort
}
