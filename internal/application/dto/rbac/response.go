package rbac

import (
	"hona/backend/internal/application/dto/address"
	"hona/backend/internal/domain/enums"
	"time"
)

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
}

type RefreshTokenResponse struct {
	AccessToken string         `json:"accessToken"`
	Roles       []RoleResponse `json:"roles"`
}

type RoleWithUsersResponse struct {
	Role  RoleResponse   `json:"role"`
	Users []UserResponse `json:"users"`
}

type AdminPetSittersListResponse struct {
	Items      []AdminPetSitterItem `json:"items"`
	Pagination PaginationResponse   `json:"pagination"`
}

type AdminPetSitterItem struct {
	ID        uint                  `json:"id"`
	FullName  string                `json:"fullName"`
	Email     string                `json:"email"`
	Phone     string                `json:"phone"`
	Status    enums.PetSitterStatus `json:"status"`
	CreatedAt time.Time             `json:"createdAt"`
	Step      enums.OnboardingStep  `json:"step"`
}

type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type WalletResponse struct {
	ID             uint    `json:"id"`
	Balance        uint    `json:"balance"`
	PendingBalance uint    `json:"pendingBalance"`
	PaymentInfo    *string `json:"paymentInfo"`
	UserID         uint    `json:"userID"`
}

type UserResponse struct {
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
	Bio             *string                      `json:"bio"`
	Wallet          WalletResponse               `json:"wallet"`
	Roles           []RoleResponse               `json:"roles"`
}
