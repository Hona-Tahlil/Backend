package seeder

import (
	"fmt"
	// "hona/backend/bootstrap"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type PetSeeder struct {
	db *gorm.DB
}

func NewPetSeeder(db *gorm.DB) *PetSeeder {
	return &PetSeeder{db: db}
}

func (s *PetSeeder) Seed(count int) error {
	if s.db.First(&entities.Pet{}).Error == nil {
		log.Println("✓ Pets already seeded, skipping...")
		return nil
	}
	// defaultProfileKey := bootstrap.Run().Env.Storage.DefaultPetProfileKey
	// if defaultProfileKey == "" {
	// 	return fmt.Errorf("default pet profile key is empty")
	// }
	var users []entities.User
	if err := s.db.Limit(100).Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}

	pets := make([]entities.Pet, 0, count)

	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]

		petKinds := enums.GetAllPetKinds()
		kind := petKinds[rand.Intn(len(petKinds))]

		var species enums.Species
		switch kind {
		case enums.Dog:
			dogSpecies := []enums.Species{
				enums.GermanShepherd,
				enums.GoldenRetriever,
				enums.Bulldog,
				enums.Poodle,
				enums.Labrador,
				enums.Husky,
			}
			species = dogSpecies[rand.Intn(len(dogSpecies))]
		case enums.Cat:
			catSpecies := []enums.Species{
				enums.Persian,
				enums.Siamese,
				enums.MaineCoon,
				enums.BritishShorthair,
			}
			species = catSpecies[rand.Intn(len(catSpecies))]
		}

		genders := []enums.PetGender{enums.MalePet, enums.FemalePet}
		gender := genders[rand.Intn(len(genders))]

		now := time.Now()
		yearsAgo := rand.Intn(16)
		daysOffset := rand.Intn(365)
		hoursOffset := rand.Intn(24)
		birthDate := now.AddDate(-yearsAgo, 0, -daysOffset).Add(-time.Duration(hoursOffset) * time.Hour)

		isAdult := yearsAgo >= 2

		var weight *float32
		var w float32
		switch kind {
		case enums.Dog:
			w = float32(10 + rand.Intn(40))
		case enums.Cat:
			w = float32(3 + rand.Intn(7))
		}
		weight = &w

		var aboutPet *string
		if rand.Intn(10) < 7 {
			about := generatePetAbout(kind)
			aboutPet = &about
		}

		name := generatePetName(kind)
		// pictureKey := defaultProfileKey

		pet := entities.Pet{
			UserID:     user.ID,
			Name:       name,
			Kind:       kind,
			Species:    species,
			BirthDate:  &birthDate,
			IsAdult:    isAdult,
			Gender:     gender,
			Weight:     weight,
			AboutPet:   aboutPet,
			// PictureKey: &pictureKey,
		}

		pets = append(pets, pet)
	}

	if err := s.db.CreateInBatches(pets, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert pets: %w", err)
	}

	log.Printf("✓ Successfully seeded %d pets", count)
	return nil
}

// generatePetName generates appropriate names based on pet kind
func generatePetName(kind enums.PetKind) string {
	dogNames := []string{
		"رکس", "ماکس", "بادی", "چارلی", "راکی", "باستر", "دیوک", "بیلی",
		"جک", "تدی", "کوپر", "بیر", "لوکی", "موچی", "برونو", "زئوس",
	}

	catNames := []string{
		"میشی", "ملوس", "تام", "گربه", "کیتی", "لونا", "سیمبا", "نالا",
		"اسکار", "میلو", "فلیکس", "گارفیلد", "پرنس", "پرنسس", "شادو", "میسی",
	}

	switch kind {
	case enums.Dog:
		return dogNames[rand.Intn(len(dogNames))]
	case enums.Cat:
		return catNames[rand.Intn(len(catNames))]
	default:
		return faker.FirstName()
	}
}

// generatePetAbout generates realistic pet descriptions
func generatePetAbout(kind enums.PetKind) string {
	dogAbouts := []string{
		"سگ بسیار مهربان و وفادار. عاشق بازی و دویدن است.",
		"انرژی بالایی دارد و نیاز به پیاده‌روی روزانه دارد.",
		"با کودکان بسیار خوب است و آرام و صبور.",
		"آموزش دیده و از دستورات پایه اطاعت می‌کند.",
		"بسیار بازیگوش و دوست‌داشتنی. عاشق توپ بازی است.",
		"نگهبان خوبی برای خانه است اما با مهمان‌ها مهربان.",
		"ترجیح می‌دهد همیشه در کنار صاحبش باشد.",
	}

	catAbouts := []string{
		"گربه آرام و مستقل. عاشق چرت زدن در آفتاب است.",
		"بسیار بازیگوش و کنجکاو. با اسباب‌بازی‌ها سرگرم می‌شود.",
		"دوست دارد نوازش شود اما سر خودش هست.",
		"شب‌ها فعال‌تر است و دوست دارد شکار بازی کند.",
		"با گربه‌های دیگر کنار می‌آید و اجتماعی است.",
		"آرام و بی‌سروصدا. برای آپارتمان مناسب است.",
		"عاشق بلندی‌هاست و دوست دارد روی قفسه‌ها بنشیند.",
	}

	switch kind {
	case enums.Dog:
		return dogAbouts[rand.Intn(len(dogAbouts))]
	case enums.Cat:
		return catAbouts[rand.Intn(len(catAbouts))]
	default:
		return "حیوان خانگی دوست‌داشتنی و مهربان."
	}
}
