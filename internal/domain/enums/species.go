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
	Beagle
	Rottweiler
	Chihuahua
	ShihTzu
	Boxer
	Pomeranian
	Dachshund
	Corgi

	// Cats
	Persian
	Siamese
	MaineCoon
	BritishShorthair
	Sphynx
	Bengal
	Ragdoll
	ScottishFold

	// Birds
	Parrot
	Canary
	Budgie
	Cockatiel
	Finch
	Lovebird
	Macaw
	Cockatoo

	// Rodents
	SyrianHamster
	DwarfHamster
	GuineaPig
	Mouse
	Rat
	Gerbil
	Chinchilla
	HollandLopRabbit
	MiniRexRabbit
	LionheadRabbit
	NetherlandDwarfRabbit

	//Other
	Other
)

func (species Species) String() string {
	switch species {

	// Dogs
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
	case Beagle:
		return "بیگل"
	case Rottweiler:
		return "روتوایلر"
	case Chihuahua:
		return "چیهواوا"
	case ShihTzu:
		return "شیه‌تزو"
	case Boxer:
		return "باکسر"
	case Pomeranian:
		return "پامرانیان"
	case Dachshund:
		return "داشهوند"
	case Corgi:
		return "کورگی"

	// Cats
	case Persian:
		return "پرشین"
	case Siamese:
		return "سیامی"
	case MaineCoon:
		return "مین کون"
	case BritishShorthair:
		return "بریتیش شورتهیر"
	case Sphynx:
		return "اسفنکس"
	case Bengal:
		return "بنگال"
	case Ragdoll:
		return "رagdoll"
	case ScottishFold:
		return "اسکاتیش فولد"

	// Birds
	case Parrot:
		return "طوطی"
	case Canary:
		return "قناری"
	case Budgie:
		return "مرغ عشق"
	case Cockatiel:
		return "کاکاتیل"
	case Finch:
		return "فنچ"
	case Lovebird:
		return "عاشق‌پرنده"
	case Macaw:
		return "ماکائو"
	case Cockatoo:
		return "کاکاتو"

	// Rodents
	case SyrianHamster:
		return "همستر سوری"
	case DwarfHamster:
		return "همستر کوتوله"
	case GuineaPig:
		return "خوکچه هندی"
	case Mouse:
		return "موش"
	case Rat:
		return "رت"
	case Gerbil:
		return "جربیل"
	case Chinchilla:
		return "چینچیلا"

	// Rabbits
	case HollandLopRabbit:
		return "خرگوش هلند لاپخ"
	case MiniRexRabbit:
		return "خرگوش مینی رِکس"
	case LionheadRabbit:
		return "خرگوش لاین‌هد"
	case NetherlandDwarfRabbit:
		return "خرگوش نِترلند دورف"

	default:
		return "دیگر نژاد ها / نمی دانم"
	}
}

func GetAllSpecies() []Species {
	return []Species{
		// Dogs
		GermanShepherd, GoldenRetriever, Bulldog, Poodle, Labrador, Husky,
		Beagle, Rottweiler, Chihuahua, ShihTzu, Boxer, Pomeranian, Dachshund, Corgi,

		// Cats
		Persian, Siamese, MaineCoon, BritishShorthair, Sphynx, Bengal, Ragdoll, ScottishFold,

		// Birds
		Parrot, Canary, Budgie, Cockatiel, Finch, Lovebird, Macaw, Cockatoo,

		// Rodents
		SyrianHamster, DwarfHamster, GuineaPig, Mouse, Rat, Gerbil, Chinchilla,
		HollandLopRabbit, MiniRexRabbit, LionheadRabbit, NetherlandDwarfRabbit,

		//Other
		Other,
	}
}

func GetSpeciesByKind(kind PetKind) []Species {
	switch kind {

	case Dog:
		return []Species{
			GermanShepherd,
			GoldenRetriever,
			Bulldog,
			Poodle,
			Labrador,
			Husky,
			Beagle,
			Rottweiler,
			Chihuahua,
			ShihTzu,
			Boxer,
			Pomeranian,
			Dachshund,
			Corgi,
		}

	case Cat:
		return []Species{
			Persian,
			Siamese,
			MaineCoon,
			BritishShorthair,
			Sphynx,
			Bengal,
			Ragdoll,
			ScottishFold,
		}

	case Bird:
		return []Species{
			Parrot,
			Canary,
			Budgie,
			Cockatiel,
			Finch,
			Lovebird,
			Macaw,
			Cockatoo,
		}

	case Rodent:
		return []Species{
			SyrianHamster,
			DwarfHamster,
			GuineaPig,
			Mouse,
			Rat,
			Gerbil,
			Chinchilla,
			HollandLopRabbit,
			MiniRexRabbit,
			LionheadRabbit,
			NetherlandDwarfRabbit,
		}

	default:
		return []Species{}
	}
}
