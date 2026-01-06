package enums
type ChatRoomStatus uint

const (
	C_Pending ChatRoomStatus = iota + 1
	C_Accepted
	C_Rejected
	C_Blocked
)
func (status ChatRoomStatus) String() string {
	switch status {
	case ChatRoomStatus(C_Pending):
		return "PENDING"
	case ChatRoomStatus(C_Accepted):
		return "ACCEPTED"
	case ChatRoomStatus(C_Rejected):
		return "REJECTED"
	case ChatRoomStatus(C_Blocked):
		return "BLOCKED"
	default:
		return "UNKNOWN"
	}
}
func GetAllChatRoomStatus() []ChatRoomStatus {
	return []ChatRoomStatus{
		C_Pending,
		C_Accepted,
		C_Rejected,
		C_Blocked,
	}
}
