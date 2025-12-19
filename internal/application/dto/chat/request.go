package chat

type CreateOrGetUserRoomRequest struct {
	PetSitterID uint 
	UserID      uint 
}


type SaveMessageRequest struct {
	RoomID  uint   
	UserID  uint   
	Content string 
}
