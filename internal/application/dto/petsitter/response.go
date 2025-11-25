package petsitter

import (
	"hona/backend/internal/domain/enums"
)

type PetSitterStatusResponse struct {
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep   `json:"onboarding_step"`
}

type PersonalInfoResponse struct {
	FirstName      string                `json:"first_name"`
	LastName       string                `json:"last_name"`
	Email          string                `json:"email"`
	Gender         enums.Gender          `json:"gender"`
	PhoneNumber    string                `json:"phone_number"`
	Province       enums.Province        `json:"province"`
	City           enums.City            `json:"city"`
	Address        string                `json:"address"`
	HouseNumber    uint                  `json:"house_number"`
	Unit           uint                  `json:"unit"`
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep  `json:"onboarding_step"`
}

type DocumentResponse struct {
	CertificateFiles []string `json:"certificate_files"`
	Files            []string `json:"files"`
}
