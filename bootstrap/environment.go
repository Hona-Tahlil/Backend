package bootstrap

import (
	"os"
	"strconv"
	"time"

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
	WebsocketSetting  WebsocketSetting
}

type Storage struct {
	Buckets               Buckets
	Endpoint              string
	AccessKey             string
	SecretKey             string
	DefaultPetProfileKey  string
	DefaultUserProfileKey string
}

type Buckets struct {
	PetProfilePic  string
	UserProfilePic string
	PetSitterFile  string
	ChatMedia      string
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
type WebsocketSetting struct {
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	PingPeriod        time.Duration
	MaxMessageSize    int
	MessageBufferSize int
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
				PetProfilePic:  os.Getenv("STORAGE_PET_PROFILE_PIC_BUCKET"),
				UserProfilePic: os.Getenv("STORAGE_USER_PROFILE_PIC_BUCKET"),
				PetSitterFile:  os.Getenv("STORAGE_PET_SITTER_FILE_BUCKET"),
				ChatMedia:      os.Getenv("STORAGE_CHAT_MEDIA_BUCKET"),
			},
			DefaultPetProfileKey:  os.Getenv("STORAGE_PET_DEFAULT_PROFILE_KEY"),
			DefaultUserProfileKey: os.Getenv("STORAGE_USER_DEFAULT_PROFILE_KEY"),
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
		WebsocketSetting: WebsocketSetting{
			WriteTimeout:      getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
			ReadTimeout:       getEnvDuration("READ_TIMEOUT", 60*time.Second),
			PingPeriod:        getEnvDuration("PING_PERIOD", 54*time.Second),
			MaxMessageSize:    getEnvInt("MAX_MESSAGE_SIZE", 524288),
			MessageBufferSize: getEnvInt("MESSAGE_BUFFER_SIZE", 256),
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

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
