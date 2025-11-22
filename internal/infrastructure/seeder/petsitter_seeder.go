package seeder

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type PetSitterSeeder struct {
	db *gorm.DB
}

func NewPetSitterSeeder(db *gorm.DB) *PetSitterSeeder {
	return &PetSitterSeeder{
		db: db,
	}
}

func (s *PetSitterSeeder) Seed(count int) error {
	petSitters := make([]entities.PetSitter, 0, count)
	var users []entities.User
	if err := s.db.Limit(100).Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}

	for i := 0; i < count; i++ {
		userID := users[i].ID

		var bio *string
		if rand.Intn(2) == 1 {
			bioText := faker.Paragraph()
			bio = &bioText
		}

		statuses := enums.GetAllPetSitterStatus()
		status := statuses[rand.Intn(len(statuses))]

		petSitter := entities.PetSitter{
			UserID: userID,
			Bio:    bio,
			Status: status,
		}

		petSitters = append(petSitters, petSitter)
	}

	if err := s.db.CreateInBatches(petSitters, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert pet sitters: %w", err)
	}

	log.Printf("✓ Successfully seeded %d pet sitters", count)
	return nil
}
