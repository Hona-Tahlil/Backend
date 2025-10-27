package enums

type BucketType uint

const (
	PetProfilePic BucketType = iota + 1
)

func (bucketType BucketType) String() string {
	switch bucketType {
	case BucketType(PetProfilePic):
		return "pet-profile-pic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		PetProfilePic,
	}
}
