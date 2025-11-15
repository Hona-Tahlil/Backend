package enums

type OnboardingStep uint8

const (
	OBS_Review  OnboardingStep = 1 
	OBS_Profile OnboardingStep = 2 
	OBS_Skills OnboardingStep = 3
	OBS_Done OnboardingStep = 4
)


func GetAllOnboardingStep() []OnboardingStep {
	return []OnboardingStep{
		OBS_Review,
		OBS_Profile,
		OBS_Skills,
		OBS_Done,
	}
}
