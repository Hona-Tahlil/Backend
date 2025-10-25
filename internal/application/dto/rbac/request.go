package rbac

type RefreshTokenRequest struct {
	RefreshToken string
}

type GetRoleByIDRequest struct {
	ID uint
}

type GetRoleByTypeRequest struct {
	Type string
}

type GetUserRolesByIDRequest struct {
	ID uint
}

type GetUserRolesByEmailRequest struct {
	Email string
}

type RemoveRoleFromUserByIDRequest struct {
	UserID uint
	RoleID uint
}

type RemoveRoleFromUserByEmailRequest struct {
	UserEmail string
	RoleID    uint
}

type AddRoleToUserByIDRequest struct {
	UserID uint
	RoleID uint
}

type AddRoleToUserByEmailRequest struct {
	UserEmail string
	RoleID    uint
}

type AddRoleRequest struct {
	Type        string
	Description *string
}

type RemoveRoleByIDRequest struct {
	ID uint
}

type RemoveRoleByTypeRequest struct {
	Type string
}

type AddPermissionToRoleRequest struct {
	RoleID       uint
	PermissionID uint
}

type RemovePermissionFromRoleRequest struct {
	RoleID       uint
	PermissionID uint
}

type GetPermissionRolesRequest struct {
	PermissionID uint
}
