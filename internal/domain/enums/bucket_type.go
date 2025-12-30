package enums

type BucketType uint

const (
	PetProfilePic BucketType = iota + 1
	UserProfilePic
	PetSitterFile
	PetSitterCert
)

func (bucketType BucketType) String() string {
	switch bucketType {
	case BucketType(PetProfilePic):
		return "pet-profile-pic"
	case BucketType(UserProfilePic):
		return "user-profile-pic"
	case BucketType(PetSitterFile):
		return "pet-sitter-file"
	case BucketType(PetSitterCert):
		return "pet-sitter-cert"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		PetProfilePic,
		UserProfilePic,
		PetSitterFile,
		PetSitterCert,
	}
}
