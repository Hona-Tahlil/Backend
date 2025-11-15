package enums

type PetSitterStatus string

const (
	PSS_Draft     PetSitterStatus = "draft"     // رکورد ایجاد شده ولی هنوز کامل نیست
	PSS_InReview  PetSitterStatus = "in_review" // ارسال برای بررسی/احراز
	PSS_Active    PetSitterStatus = "active"    // قابل نمایش و رزرو
	PSS_Rejected  PetSitterStatus = "rejected"  // رد شده (نیاز به اصلاح)
	PSS_Suspended PetSitterStatus = "suspended" // تعلیق موقت/دائمی
)



func GetAllPetSitterStatus() []PetSitterStatus {
	return []PetSitterStatus{
		PSS_Draft,
		PSS_Active,
		PSS_InReview,
		PSS_Rejected,
		PSS_Suspended,
	}
}
