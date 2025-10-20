package rbac

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RefreshTokenResponse struct {
	AccessToken string               `json:"accessToken"`
	Permissions []PermissionResponse `json:"permissions"`
}
