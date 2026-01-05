package petsitter

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/dto/pet"
	"hona/backend/internal/application/dto/rbac"
	"hona/backend/internal/application/dto/servicedto"
	"hona/backend/internal/domain/enums"
	"time"
)

type PetSitterStatusResponse struct {
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep  `json:"onboarding_step"`
}

type PersonalInfoResponse struct {
	FirstName      string                `json:"first_name"`
	LastName       string                `json:"last_name"`
	Email          string                `json:"email"`
	Gender         enums.Gender          `json:"gender"`
	BirthDate      *time.Time            `json:"birth_date"`
	PhoneNumber    string                `json:"phone_number"`
	Province       enums.Province        `json:"province"`
	City           enums.City            `json:"city"`
	Address        string                `json:"address"`
	HouseNumber    uint                  `json:"house_number"`
	Unit           uint                  `json:"unit"`
	UserProfilePic *string               `json:"user_profile_pic"`
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep  `json:"onboarding_step"`
}

type DocumentResponse struct {
	CertificateFiles []string `json:"certificate_files"`
	Files            []string `json:"files"`
}

type PetSitterListItemResponse struct {
	ID             uint                  `json:"id"`
	UserID         uint                  `json:"user_id"`
	FirstName      string                `json:"first_name"`
	LastName       string                `json:"last_name"`
	Email          string                `json:"email"`
	PhoneNumber    string                `json:"phone_number"`
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep  `json:"onboarding_step"`
	CreatedAt      string                `json:"created_at"`
}

type PetSittersListResponse struct {
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Count      int                         `json:"count"`
	PetSitters []PetSitterListItemResponse `json:"pet_sitters"`
}

type PetSitterDetailsResponse struct {
	PersonalInfo   rbac.UserResponse     `json:"personalInfo"`
	Documents      DocumentResponse      `json:"documents"`
	Skills         SkillsResponse        `json:"skills"`
	Status         enums.PetSitterStatus `json:"status"`
	OnboardingStep enums.OnboardingStep  `json:"onboarding_step"`
	CreatedAt      string                `json:"created_at"`
}

type PetSitterSelfProfileResponse struct {
	ID              uint                             `json:"id"`
	PetSitterID     uint                             `json:"petSitterID"`
	Email           string                           `json:"email"`
	IsEmailVerified bool                             `json:"isEmailVerified"`
	FirstName       string                           `json:"firstName"`
	LastName        string                           `json:"lastName"`
	Address         *address.AddressInfoResponse     `json:"address"`
	Phone           *string                          `json:"phone"`
	IsPhoneVerified bool                             `json:"isPhoneVerified"`
	Gender          string                           `json:"gender"`
	BirthDate       *time.Time                       `json:"birthDate"`
	PictureLink     *string                          `json:"pictureLink"`
	Bio             *string                          `json:"bio"`
	Services        []servicedto.ServiceInfoResponse `json:"services"`
	PetKinds        []pet.PetKindResponse            `json:"petKinds"`
}

type PetSitterProfileResponse struct {
	ID          uint                             `json:"id"`
	UserID      uint                             `json:"user_id"`
	FirstName   string                           `json:"first_name"`
	LastName    string                           `json:"last_name"`
	PictureLink *string                          `json:"picture_link"`
	Province    string                           `json:"province"`
	City        string                           `json:"city"`
	Bio         string                           `json:"bio"`
	Services    []servicedto.ServiceInfoResponse `json:"services"`
	PetKinds    []pet.PetKindResponse            `json:"pet_kinds"`
	CreatedAt   string                           `json:"created_at"`
}

type SkillsResponse struct {
	Bio      string                           `json:"bio"`
	PetKinds []pet.PetKindResponse            `json:"pet_kinds"`
	Services []servicedto.ServiceInfoResponse `json:"services"`
}

type SearchPetSitterInfoResponse struct {
	ID        uint     `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Province  string   `json:"province"`
	City      string   `json:"city"`
	Services  []string `json:"services"`
	MinPrice  uint     `json:"min_price"`
	PetKinds  []string `json:"pet_kinds"`
	Rate      float64  `json:"rate"`
	Comments  uint     `json:"comments"`
}
