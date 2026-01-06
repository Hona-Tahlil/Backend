package bootstrap

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
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
	Logger            LoggerConfig
	WebsocketSetting  WebsocketSetting
	RabbitMQ          RabbitMQ
	RateLimit         RateLimit
}

type RateLimit struct {
	Limit int
	Burst int
}

type RabbitMQ struct {
	User          string
	Password      string
	Host          string
	Port          string
	VHost         string
	MaxRetryCount int
	RetryDelay    time.Duration
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
	PetSitterCert  string
	PetSitterFile  string
	UserProfilePic string
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

type LoggerConfig struct {
	Level        slog.Level
	TextStdout   bool
	JSONStdout   bool
	JSONFilePath string
	AddSource    bool
	ServiceName  string
	Environment  string
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
				PetSitterCert:  os.Getenv("STORAGE_PET_SITTER_CERT_BUCKET"),
				PetSitterFile:  os.Getenv("STORAGE_PET_SITTER_FILE_BUCKET"),
				UserProfilePic: os.Getenv("STORAGE_USER_PROFILE_PIC_BUCKET"),
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
		Logger: loadLoggerConfig(),
		WebsocketSetting: WebsocketSetting{
			WriteTimeout:      getEnvDuration("WRITE_TIMEOUT", 10*time.Second),
			ReadTimeout:       getEnvDuration("READ_TIMEOUT", 60*time.Second),
			PingPeriod:        getEnvDuration("PING_PERIOD", 54*time.Second),
			MaxMessageSize:    getEnvInt("MAX_MESSAGE_SIZE", 524288),
			MessageBufferSize: getEnvInt("MESSAGE_BUFFER_SIZE", 256),
		},
		RabbitMQ: RabbitMQ{
			User:          os.Getenv("AMQP_USER"),
			Password:      os.Getenv("AMQP_PASSWORD"),
			Host:          os.Getenv("AMQP_HOST"),
			Port:          os.Getenv("AMQP_PORT"),
			VHost:         os.Getenv("AMQP_VHOST"),
			MaxRetryCount: getEnvInt("AMQP_MAX_RETRY", 3),
			RetryDelay:    getEnvDuration("AMQP_RETRY_DELAY", 5*time.Second),
		},
		RateLimit: RateLimit{
			Limit: getEnvInt("RATE_LIMIT", 5),
			Burst: getEnvInt("RATE_lIMIT_BURST", 10),
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

func loadLoggerConfig() LoggerConfig {
	const (
		defaultLogFile    = "logs/app.json"
		defaultService    = "backend"
		defaultEnv        = "development"
		defaultTextStdout = true
	)

	level := parseLogLevel(os.Getenv("LOG_LEVEL"))
	textStdout := parseBool(os.Getenv("LOG_TEXT_STDOUT"), defaultTextStdout)
	jsonStdout := parseBool(os.Getenv("LOG_JSON_STDOUT"), false)
	jsonFilePath := os.Getenv("LOG_FILE_PATH")
	if jsonFilePath == "" && parseBool(os.Getenv("LOG_JSON_FILE"), true) {
		jsonFilePath = defaultLogFile
	}
	if !parseBool(os.Getenv("LOG_JSON_FILE"), true) {
		jsonFilePath = ""
	}
	addSource := parseBool(os.Getenv("LOG_ADD_SOURCE"), false)
	serviceName := strings.TrimSpace(os.Getenv("LOG_SERVICE_NAME"))
	if serviceName == "" {
		serviceName = defaultService
	}
	environment := strings.TrimSpace(os.Getenv("LOG_ENV"))
	if environment == "" {
		environment = strings.TrimSpace(os.Getenv("APP_ENV"))
	}
	if environment == "" {
		environment = defaultEnv
	}

	return LoggerConfig{
		Level:        level,
		TextStdout:   textStdout,
		JSONStdout:   jsonStdout,
		JSONFilePath: jsonFilePath,
		AddSource:    addSource,
		ServiceName:  serviceName,
		Environment:  environment,
	}
}

func parseBool(value string, defaultVal bool) bool {
	if value == "" {
		return defaultVal
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return defaultVal
	}
}

func parseLogLevel(value string) slog.Level {
	if value == "" {
		return slog.LevelInfo
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		if parsed, err := strconv.Atoi(value); err == nil {
			return slog.Level(parsed)
		}
		return slog.LevelInfo
	}
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}
