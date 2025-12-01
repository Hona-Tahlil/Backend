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

		db.AutoMigrate(
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
		)
	})

	return dbInstance
}
