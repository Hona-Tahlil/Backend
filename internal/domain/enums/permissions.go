package enums

type Permission uint

const (
	RequestPermission Permission = iota + 1
)

func (permission Permission) String() string {
	switch permission {
	case Permission(RequestPermission):
		return "request"
	}
	return ""
}

func GetAllPermissionTypes() []Permission {
	return []Permission{
		RequestPermission,
	}
}
