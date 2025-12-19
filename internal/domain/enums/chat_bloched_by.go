package enums
type ChatBlockedBy uint

const (
	C_User ChatBlockedBy = iota + 1
	C_PetSitter
)
func (blockedBy ChatBlockedBy) String() string {
	switch blockedBy {
	case ChatBlockedBy(C_User):
		return "USER"
	case ChatBlockedBy(C_PetSitter):
		return "PET_SITTER"
	default:
		return "UNKNOWN"
	}
}
func GetAllChatBlockedBy() []ChatBlockedBy {
	return []ChatBlockedBy{
		C_User,
		C_PetSitter,
	}
}