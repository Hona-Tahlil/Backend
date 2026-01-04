package seeder

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct {
	db *gorm.DB
}

func NewUserSeeder(db *gorm.DB) *UserSeeder {
	return &UserSeeder{db: db}
}

func (s *UserSeeder) Seed(count int) error {
	users := make([]entities.User, 0, count)

	var userCount int64
	s.db.Model(&entities.User{}).Count(&userCount)
	log.Println("Current user count in database:", userCount)
	if userCount > 0 {
		log.Println("✓ Users already seeded, skipping...")
		return nil
	}
	// defaultProfileKey := bootstrap.Run().Env.Storage.DefaultUserProfileKey
	// if defaultProfileKey == "" {
	// 	return fmt.Errorf("default user profile key is empty")
	// }

	for i := 0; i < count; i++ {
		email := "test" + fmt.Sprintf("%d", i) + "@email.com"
		firstName := faker.FirstName()
		lastName := faker.LastName()

		rawPass := "password123"
		hashed, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}
		log.Println("Generated user:", email)
		var phone *string
		if rand.Intn(2) == 1 {
			p := faker.Phonenumber()
			phone = &p
		}

		var birthDate *time.Time
		if rand.Intn(2) == 1 {
			now := time.Now()
			yearsAgo := 18 + rand.Intn(62) // 18 to 80 years
			daysOffset := rand.Intn(365)
			t := now.AddDate(-yearsAgo, 0, -daysOffset)
			birthDate = &t
		}

		genders := []enums.Gender{enums.Male, enums.Female}
		gender := genders[rand.Intn(len(genders))]
		// pictureKey := defaultProfileKey

		user := entities.User{
			Email:           email,
			IsEmailVerified: i%2 == 0,
			Password:        string(hashed),
			FirstName:       firstName,
			LastName:        lastName,
			Phone:           phone,
			IsPhoneVerified: i%3 == 0,
			Gender:          gender,
			BirthDate:       birthDate,
			// PictureLink:     &pictureKey,
		}

		users = append(users, user)
	}
	log.Println("Seeding users...")
	if err := s.db.CreateInBatches(users, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert users: %w", err)
	}

	log.Printf("✓ Successfully seeded %d users", count)
	return nil
}
