package enums

type PetSitterStatus uint

const (
	PSS_Draft    PetSitterStatus = iota + 1
	PSS_InReview  
	PSS_Active   
	PSS_Rejected  
	PSS_Suspended 
)

func (petSitterStatus PetSitterStatus) String() string {
	switch petSitterStatus {
	case PetSitterStatus(PSS_Draft):
		return "draft"
	case PetSitterStatus(PSS_InReview):
		return "in_review"
	case PetSitterStatus(PSS_Active):
		return "active"
	case PetSitterStatus(PSS_Rejected):
		return "rejected"
	case PetSitterStatus(PSS_Suspended):
		return "suspended"
	}
	return ""
}

func GetAllPetSitterStatus() []PetSitterStatus {
	return []PetSitterStatus{
		PSS_Draft,
		PSS_Active,
		PSS_InReview,
		PSS_Rejected,
		PSS_Suspended,
	}
}
