package petsitter

import "hona/backend/internal/domain/enums"

type PetSitterStatusResponse struct {
	UserID         uint
	Status         enums.PetSitterStatus
	OnboardingStep enums.OnboardingStep
}
