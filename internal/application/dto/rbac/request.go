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

type GetRoleWithUsersByTypeRequest struct {
	Type  string
	Page  int
	Count int
}

type GetRoleWithUsersByIDRequest struct {
	ID    uint
	Page  int
	Count int
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

type ListRolesWithUsersRequest struct {
	Page  int
	Count int
}


type ListPetSittersRequest struct {
    Page    int
    Limit   int
    Status  *string // optional
    Search  *string // optional
    Sort    string  // created_at_desc (default)
}