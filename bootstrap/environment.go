package bootstrap

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DSN string
}

func NewEnv() *Env {
	godotenv.Load(".env")
	return &Env{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		),
	}
}
