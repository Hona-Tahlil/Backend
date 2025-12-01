package enums

type OnboardingStep uint8

const (
<<<<<<< HEAD
	OBS_Review  OnboardingStep = iota + 1
	OBS_Profile 
	OBS_Documents
	OBS_Done 
)


=======
	OBS_Review OnboardingStep = iota + 1
	OBS_Profile
	OBS_Documents
	OBS_Done
)

>>>>>>> dev
func GetAllOnboardingStep() []OnboardingStep {
	return []OnboardingStep{
		OBS_Review,
		OBS_Profile,
		OBS_Documents,
		OBS_Done,
	}
}
