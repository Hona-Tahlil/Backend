package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PrimaryDB    Database
	PrimaryRedis Redis
}

type Database struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type Redis struct {
	Port      string
	Address   string
	Password  string
	RDBNumber string
}

func NewEnv() *Env {
	godotenv.Load(".env")
	return &Env{
		PrimaryDB: Database{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		PrimaryRedis: Redis{
			Port:      os.Getenv("RDB_PORT"),
			Address:   os.Getenv("RDB_ADDRESS"),
			Password:  os.Getenv("RDB_PASSWORD"),
			RDBNumber: os.Getenv("RDB_NUMBER"),
		},
	}
}
