package enums
type ChatBlockedBy uint

const (
	C_BlockedByUser ChatBlockedBy = iota + 1
	C_BlockedByPetSitter
	C_BlockedByAdmin
)
func (blockedBy ChatBlockedBy) String() string {
	switch blockedBy {
	case ChatBlockedBy(C_BlockedByUser):
		return "USER"
	case ChatBlockedBy(C_BlockedByPetSitter):
		return "PET_SITTER"
	case ChatBlockedBy(C_BlockedByAdmin):
		return "ADMIN"		
	default:
		return "UNKNOWN"
	}
}
func GetAllChatBlockedBy() []ChatBlockedBy {
	return []ChatBlockedBy{
		C_BlockedByUser,
		C_BlockedByPetSitter,
		C_BlockedByAdmin,
	}
}