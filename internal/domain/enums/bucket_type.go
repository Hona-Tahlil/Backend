package enums

type BucketType uint

const (
	PetProfilePic BucketType = iota + 1
<<<<<<< HEAD
	UserProfilePic
	PetSitterFile
=======
>>>>>>> dev
)

func (bucketType BucketType) String() string {
	switch bucketType {
	case BucketType(PetProfilePic):
		return "pet-profile-pic"
<<<<<<< HEAD
	case BucketType(UserProfilePic):
		return "user-profile-pic"
	case BucketType(PetSitterFile):
		return "pet-sitter-file"
=======
>>>>>>> dev
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		PetProfilePic,
<<<<<<< HEAD
		UserProfilePic,
		PetSitterFile,
=======
>>>>>>> dev
	}
}
