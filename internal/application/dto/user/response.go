package user

import (
	"hona/backend/internal/application/dto/rbac"
)

type LoginResponse struct {
	AccessToken string              `json:"accessToken"`
	Roles       []rbac.RoleResponse `json:"roles"`
}

type MLData struct {
	Token string `json:"ml"`
}
