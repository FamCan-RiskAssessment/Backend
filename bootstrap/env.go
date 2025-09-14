package bootstrap

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Server     Server
	Database   Database
	Cache      Redis
	SMSGateway SMSGateway
	Pagination Pagination
	OTP        OTP
	SuperAdmin SuperAdmin
}

type Server struct {
	Port string
	Mode string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Redis struct {
	Address   string
	Port      string
	Password  string
	RDBNumber string
}

type SMSGateway struct {
	APIKey string
}

type Pagination struct {
	DefaultPage     int
	DefaultPageSize int
}

type OTP struct {
	Length       int
	ExpiryMinute int
	MaxAttempts  int
}

type SuperAdmin struct {
	Phone    string
	Password string
}

func NewEnv() *Env {
	godotenv.Load(".env")
	return &Env{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
			Mode: os.Getenv("SERVER_MODE"),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Cache: Redis{
			Port:      os.Getenv("RDB_PORT"),
			Address:   os.Getenv("RDB_ADDRESS"),
			Password:  os.Getenv("RDB_PASSWORD"),
			RDBNumber: os.Getenv("RDB_NUMBER"),
		},
		SMSGateway: SMSGateway{
			APIKey: os.Getenv("SMS_API_KEY"),
		},
		Pagination: Pagination{
			DefaultPage:     getEnvInt("PAGINATION_DEFAULT_PAGE", 1),
			DefaultPageSize: getEnvInt("PAGINATION_DEFAULT_PAGE_SIZE", 10),
		},
		OTP: OTP{
			Length:       getEnvInt("OTP_LENGTH", 6),
			ExpiryMinute: getEnvInt("OTP_EXPIRY_MINUTES", 2),
			MaxAttempts:  getEnvInt("OTP_MAX_ATTEMPTS", 3),
		},
		SuperAdmin: SuperAdmin{
			Phone:    os.Getenv("SUPER_ADMIN_PHONE"),
			Password: os.Getenv("SUPER_ADMIN_PASSWORD"),
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
