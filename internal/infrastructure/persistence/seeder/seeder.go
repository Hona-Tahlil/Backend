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

func SeedUsers(db *gorm.DB, count int) error {
	for i := 0; i < count; i++ {
		// random base data
		email := faker.Email()
		firstName := faker.FirstName()
		lastName := faker.LastName()

		rawPass := faker.Password()
		hashed, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		var phone *string
		if rand.Intn(2) == 1 { // 50% chance phone exists
			p := faker.Phonenumber()
			phone = &p
		}

		var birthDate *time.Time
		if rand.Intn(2) == 1 {
			// t := faker.Date()
			// birthDate = &t
		}

		// random gender
		gender := enums.Gender(rand.Intn(2)) // assuming 0=Male,1=Female

		user := entities.User{
			Email:           email,
			IsEmailVerified: rand.Intn(2) == 1,
			Password:        string(hashed),
			FirstName:       firstName,
			LastName:        lastName,
			Phone:           phone,
			IsPhoneVerified: rand.Intn(2) == 1,
			Gender:          gender,
			BirthDate:       birthDate,
			PictureLink:     nil,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("failed to seed user %s: %v\n", email, err)
			return err
		}

		log.Printf("Seeded user: %s (password: %s)", email, rawPass)
	}
	return nil
}
