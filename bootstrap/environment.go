package bootstrap

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PrimaryDB         Database
	TokenExpires      TokenExpires
	Storage           Storage
	PrimaryRedis      Redis
	EmailConfig       EmailConfig
	URLs              URLs
	EmailVerification EmailVerification
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

type EmailVerification struct {
	ExpireMinutes int
}

type URLs struct {
	BaseURL string
}

type Redis struct {
	Port      string
	Address   string
	Password  string
	RDBNumber string
}
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewEnv() *Env {
	godotenv.Load(".env")
	expireMinutes, _ := strconv.Atoi(os.Getenv("EMAIL_EXPIRE_MINUTES"))
	return &Env{
		PrimaryDB: Database{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		TokenExpires: TokenExpires{
			LongRefreshHours:  getEnvInt("LONG_REFRESH_HOURS", 168),
			ShortRefreshHours: getEnvInt("SHORT_REFRESH_HOURS", 48),
			AccessMinutes:     getEnvInt("ACCESS_MINUTES", 500),
		},
		Storage: Storage{
			Endpoint:  os.Getenv("STORAGE_ENDPOINT"),
			AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
			SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
			Buckets: Buckets{
				PetProfilePic: os.Getenv("STORAGE_PET_PROFILE_PIC_BUCKET"),
			},
		},
		PrimaryRedis: Redis{
			Port:      os.Getenv("RDB_PORT"),
			Address:   os.Getenv("RDB_ADDRESS"),
			Password:  os.Getenv("RDB_PASSWORD"),
			RDBNumber: os.Getenv("RDB_NUMBER"),
		},
		EmailConfig: EmailConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
		URLs: URLs{
			BaseURL: os.Getenv("BASE_URL"),
		},
		EmailVerification: EmailVerification{
			ExpireMinutes: expireMinutes,
		},
	}
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
