package enums

type OnboardingStep uint8

const (
	OBS_Review  OnboardingStep = iota + 1
	OBS_Profile 
	OBS_Documents
	OBS_Done 
)


func GetAllOnboardingStep() []OnboardingStep {
	return []OnboardingStep{
		OBS_Review,
		OBS_Profile,
		OBS_Documents,
		OBS_Done,
	}
}
