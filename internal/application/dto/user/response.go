package user

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/application/dto/rbac"
	"time"
)

type LoginResponse struct {
	AccessToken string              `json:"accessToken"`
	Roles       []rbac.RoleResponse `json:"roles"`
}

type MLData struct {
	Token string `json:"ml"`
}

type ProfileResponse struct {
	ID              uint                         `json:"id"`
	Email           string                       `json:"email"`
	IsEmailVerified bool                         `json:"isEmailVerified"`
	FirstName       string                       `json:"firstName"`
	LastName        string                       `json:"lastName"`
	Address         *address.AddressInfoResponse `json:"address"`
	Phone           *string                      `json:"phone"`
	IsPhoneVerified bool                         `json:"isPhoneVerified"`
	Gender          string                       `json:"gender"`
	BirthDate       *time.Time                   `json:"birthDate"`
	PictureLink     *string                      `json:"pictureLink"`
}
