package enums

type Role uint

const (
	GuestRole Role = iota + 1
	PetOwnerRole
	PetSitterRole
	PetOwnerSitterRole
	AdminRole
)

func (role Role) String() string {
	switch role {
	case GuestRole:
		return "guest"
	case PetOwnerRole:
		return "owner"
	case PetSitterRole:
		return "sitter"
	case AdminRole:
		return "admin"
	case PetOwnerSitterRole:
		return "owner/sitter"
	}
	return ""
}

func GetAllRoleTypes() []Role {
	return []Role{
		GuestRole,
		PetOwnerRole,
		PetSitterRole,
		PetOwnerSitterRole,
		AdminRole,
	}
}
