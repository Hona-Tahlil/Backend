package rbac

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
