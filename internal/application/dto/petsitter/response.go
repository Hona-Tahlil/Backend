package petsitter

import (
	"hona/backend/internal/domain/enums"
	"mime/multipart"
	"time"
)

type PetSitterStatusResponse struct {
	UserID         uint
	Status         enums.PetSitterStatus
	OnboardingStep enums.OnboardingStep
}

type PersonalInfoResponse struct {
	UserID         uint
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	Gender         enums.Gender   `json:"gender,omitempty"`
	BirthDate      *time.Time     `json:"birth_date,omitempty"`
	PhoneNumber    string         `json:"phone_number,omitempty"`
	Province       enums.Province `json:"province" validate:"required"`
	City           enums.City     `json:"city,omitempty"`
	Address        string         `json:"address,omitempty"`
	Pelak          uint           `json:"pelak" validate:"required"`
	Vahed          uint           `json:"vahed" validate:"required"`
	PostalCode     string         `json:"postalCode" validate:"required,len=10"`
	Status         enums.PetSitterStatus
	OnboardingStep enums.OnboardingStep
}

type DocumentResponse struct {
	UserID          uint
	CertificateFile []*multipart.FileHeader
	File            *multipart.FileHeader
	Status          enums.PetSitterStatus
	OnboardingStep  enums.OnboardingStep
}
