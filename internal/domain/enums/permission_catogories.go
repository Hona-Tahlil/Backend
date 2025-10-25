package enums

type PermissionCategory uint

const (
	ReadPermissionCategory PermissionCategory = iota + 1
)

func (permissionCategory PermissionCategory) String() string {
	switch permissionCategory {
	case PermissionCategory(ReadPermissionCategory):
		return "request"
	}
	return ""
}

func GetAllPermissionCategoryTypes() []PermissionCategory {
	return []PermissionCategory{
		ReadPermissionCategory,
	}
}
