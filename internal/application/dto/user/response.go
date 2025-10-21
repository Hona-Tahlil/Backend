package user

import "hona/backend/internal/application/dto/rbac"

type LoginResponse struct {
	AccessToken string                    `json:"accessToken"`
	IsVerified  bool                      `json:"isVerified"`
	Permissions []rbac.PermissionResponse `json:"permissions"`
}
