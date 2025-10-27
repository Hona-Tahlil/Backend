package enums

type Species uint

const (
	// Dogs
	GermanShepherd Species = iota + 1
	GoldenRetriever
	Bulldog
	Poodle
	Labrador
	Husky

	// Cats
	Persian
	Siamese
	MaineCoon
	BritishShorthair

	// Birds
	Parrot
	Canary
	Budgie
	Cockatiel

	// Others
	Goldfish
	Hamster
)

func (species Species) String() string {
	switch species {
	case GermanShepherd:
		return "ژرمن شپرد"
	case GoldenRetriever:
		return "گلدن رتریور"
	case Bulldog:
		return "بولداگ"
	case Poodle:
		return "پودل"
	case Labrador:
		return "لابرادور"
	case Husky:
		return "هاسکی"
	case Persian:
		return "پرشین"
	case Siamese:
		return "سیامی"
	case MaineCoon:
		return "مین کون"
	case BritishShorthair:
		return "بریتیش شورتهیر"
	case Parrot:
		return "طوطی"
	case Canary:
		return "قناری"
	case Budgie:
		return "مرغ عشق"
	case Cockatiel:
		return "کاکاتیل"
	case Goldfish:
		return "ماهی قرمز"
	case Hamster:
		return "همستر"
	default:
		return ""
	}
}

func GetAllSpecies() []Species {
	return []Species{
		GermanShepherd, GoldenRetriever, Bulldog, Poodle, Labrador, Husky,
		Persian, Siamese, MaineCoon, BritishShorthair,
		Parrot, Canary, Budgie, Cockatiel,
		Goldfish, Hamster,
	}
}
