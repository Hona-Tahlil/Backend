package enums

type BucketType uint

const (
	PetProfilePic BucketType = iota + 1
	UserProfilePic
	PetSitterFile
	ChatMedia
)

func (bucketType BucketType) String() string {
	switch bucketType {
	case BucketType(PetProfilePic):
		return "pet-profile-pic"
	case BucketType(UserProfilePic):
		return "user-profile-pic"
	case BucketType(PetSitterFile):
		return "pet-sitter-file"
	case BucketType(ChatMedia):
		return "chat-media"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		PetProfilePic,
		UserProfilePic,
		PetSitterFile,
		ChatMedia,
	}
}
