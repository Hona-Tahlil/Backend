package bootstrap

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PrimaryDB    Database
	TokenExpires TokenExpires
	Storage      Storage
}

type Storage struct {
	Buckets   Buckets
	Endpoint  string
	AccessKey string
	SecretKey string
}

type Buckets struct {
	PetProfilePic string
}

type TokenExpires struct {
	LongRefreshHours  int
	ShortRefreshHours int
	AccessMinutes     int
}

type Database struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

func NewEnv() *Env {
	LongRefreshHours, _ := strconv.Atoi(os.Getenv("LONG_REFRESH_HOURS"))
	ShortRefreshHours, _ := strconv.Atoi(os.Getenv("SHORT_REFRESH_HOURS"))
	AccessMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_MINUTES"))
	godotenv.Load(".env")
	return &Env{
		PrimaryDB: Database{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		TokenExpires: TokenExpires{
			LongRefreshHours:  LongRefreshHours,
			ShortRefreshHours: ShortRefreshHours,
			AccessMinutes:     AccessMinutes,
		},
		Storage: Storage{
			Endpoint:  os.Getenv("STORAGE_ENDPOINT"),
			AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
			SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
			Buckets: Buckets{
				PetProfilePic: os.Getenv("STORAGE_PET_PROFILE_PIC_BUCKET"),
			},
		},
	}
}
