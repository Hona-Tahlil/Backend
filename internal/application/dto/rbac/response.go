package rbac

import (
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
	Role  RoleResponse       `json:"role"`
	Users []UserInfoResponse `json:"users"`
}

type UserInfoResponse struct {
	Email string `json:"email"`
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
