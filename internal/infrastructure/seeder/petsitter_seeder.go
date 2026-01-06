package seeder

import (
	"fmt"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"
	"log/slog"

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
	if s.db.First(&entities.PetSitter{}).Error == nil {
		slog.Info("✓ Pet sitters already seeded, skipping...")
		return nil
	}

	var users []entities.User
	if err := s.db.Limit(1000).Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}
	if len(users) == 0 {
		return fmt.Errorf("no users found, please seed users first")
	}
	if len(users) < count {
		count = len(users)
	}

	petSitters := make([]entities.PetSitter, 0, count)

	for i := 0; i < count; i++ {
		userID := users[i].ID

		bioText := faker.Paragraph()
		status := enums.PSS_Draft
		petSitter := entities.PetSitter{
			UserID:       userID,
			Bio:          &bioText,
			Status:       status,
			OnboardingStep: enums.OBS_Review,
			PetKinds: randomPetKinds(), 
		}

		petSitters = append(petSitters, petSitter)
	}

	if err := s.db.CreateInBatches(petSitters, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert pet sitters: %w", err)
	}

	slog.Info("✓ Successfully seeded pet sitters", "count", count)
	return nil
}
