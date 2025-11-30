package persistence

import (
	"fmt"
	"log"
	"os"
	"sync"

	"hona/backend/bootstrap"
	"hona/backend/internal/domain/entities"
	"hona/backend/internal/domain/enums"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var dbOnce sync.Once

type dbConfigStruct struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func NewPostgresDatabase() *gorm.DB {
	raw := bootstrap.Run().Env.PrimaryDB

	cfg := dbConfigStruct{
		Host:     raw.Host,
		User:     raw.User,
		Password: raw.Password,
		Name:     raw.Name,
		Port:     raw.Port,
	}

	if cfg.Host == "" {
		cfg.Host = os.Getenv("DB_HOST")
	}
	if cfg.User == "" {
		cfg.User = os.Getenv("DB_USER")
	}
	if cfg.Password == "" {
		cfg.Password = os.Getenv("DB_PASSWORD")
	}
	if cfg.Name == "" {
		cfg.Name = os.Getenv("DB_NAME")
	}
	if cfg.Port == "" {
		cfg.Port = os.Getenv("DB_PORT")
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Name == "" || cfg.Port == "" {
		log.Fatalf("[DB CONFIG ERROR] Incomplete DB config after env fallback: host=%q user=%q dbname=%q port=%q",
			cfg.Host, cfg.User, cfg.Name, cfg.Port,
		)
	}

	fmt.Printf(
		"[DB] Connecting → host=%s user=%s dbname=%s port=%s\n",
		cfg.Host,
		cfg.User,
		cfg.Name,
		cfg.Port,
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)

	dbOnce.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("❌ Failed to connect to Postgres: %v", err)
		}

		dbInstance = db

		if err = db.Migrator().DropTable(
			&entities.User{},
			&entities.Role{},
			&entities.Permission{},
			&entities.Wallet{},
			&entities.Request{},
			&entities.CalendarSlot{},
			&entities.Pet{},
			&entities.PetSitter{},
			&entities.Service{},
			&entities.Chat{},
			&entities.Comment{},
			&entities.Address{},
			&entities.Province{},
			&entities.City{},
			&entities.TextMessage{},
			&entities.Transaction{},
			&entities.Transfer{},
		); err != nil {
			log.Fatalf("❌ Failed to drop tables: %v", err)
		}

		if err := db.AutoMigrate(
			&entities.User{},
			&entities.Role{},
			&entities.Permission{},
			&entities.Wallet{},
			&entities.Request{},
			&entities.CalendarSlot{},
			&entities.Pet{},
			&entities.PetSitter{},
			&entities.Service{},
			&entities.Chat{},
			&entities.Comment{},
			&entities.Address{},
			&entities.Province{},
			&entities.City{},
			&entities.TextMessage{},
			&entities.Transaction{},
			&entities.Transfer{},
		); err != nil {
			log.Fatalf("❌ AutoMigrate failed: %v", err)
		}

		pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := entities.User{
			Email:           "test1@email.com",
			Password:        string(pass),
			FirstName:       "John",
			LastName:        "Doe",
			IsEmailVerified: true,
			Gender:          1,
		}
		db.Create(&user)
		var foundUser entities.User
		err = db.First(&foundUser, "email = ?", "test1@email.com").Error
		province := entities.Province{
			Name: 1,
		}
		db.Create(&province)
		var foundProvince entities.Province
		err = db.First(&foundProvince, "name = ?", 1).Error
		// err := db.First(&foundUser, "email = ?", "test1@email.com").Error
		// var foundUser entities.User
		city := entities.City{
			Name:       11,
			ProvinceID: foundProvince.ID,
		}
		db.Create(&city)
		var foundCity entities.City
		err = db.First(&foundCity, "name = ?", 11).Error
		address := entities.Address{
			Province:      foundProvince,
			ProvinceID:    foundProvince.ID,
			City:          foundCity,
			CityID:        foundCity.ID,
			StreetAddress: "123 Main St",
			HouseNumber:   123,
			Unit:          1,
			Refer: 	  foundUser.ID,
			Type:    "User",
			// OwnerID:       foundUser.ID,
		}
		db.Create(&address)
		var foundAddress entities.Address
		err = db.First(&foundAddress, "refer = ?", foundUser.ID).Error
		foundUser.Address = &foundAddress
		db.Save(foundUser)
		petsitter := entities.PetSitter{
			UserID:      foundUser.ID,
			Status: enums.PSS_Draft,
			OnboardingStep: enums.OBS_Review,
		}
		db.Create(&petsitter)
		fmt.Println("✅ Database initialized successfully.")
	})

	return dbInstance
}
