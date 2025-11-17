package enums

type BucketType uint

const (
	PetProfilePic BucketType = iota + 1
	PetSitterFile
)

func (bucketType BucketType) String() string {
	switch bucketType {
	case BucketType(PetProfilePic):
		return "pet-profile-pic"
	case BucketType(PetSitterFile):
		return "pet-sitter-file"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		PetProfilePic,
		PetSitterFile,
	}
}
