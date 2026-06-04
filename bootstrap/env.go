package bootstrap

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	Server          Server
	Database        Database
	Cache           Redis
	SMSGateway      SMSGateway
	Pagination      Pagination
	OTP             OTP
	SuperAdmin      SuperAdmin
	MinIO           MinIO
	CalcURL         CalcURL
	VerificationAPI VerificationAPI
	Security        Security
	JWT             JWT
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
	Provider string
	APIKey   string
	Username string
	Password string
	Source   string
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

type MinIO struct {
	Bucket    string
	Prefixes  BucketPrefix
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type BucketPrefix struct {
	Mamography        string
	Cancer            string
	GeneticTest       string
	FatherGeneticTest string
	MotherGeneticTest string
}

type SuperAdmin struct {
	Phone    string
	Password string
}

type CalcURL struct {
	Premm5 string
	BCRA   string
	Gail   string
	PLCO   string
}

type VerificationAPI struct {
	BaseURL string
	APIKey  string
}

type Security struct {
	EncryptionKey      string
	RateLimitPerMinute int
	RateLimitWindow    int
}

type JWT struct {
	PrivateKey string
	PublicKey  string
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
			Provider: os.Getenv("SMS_PROVIDER"),
			APIKey:   os.Getenv("SMS_API_KEY"),
			Username: os.Getenv("SMS_USERNAME"),
			Password: os.Getenv("SMS_PASSWORD"),
			Source:   os.Getenv("SMS_SOURCE"),
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
		MinIO: MinIO{
			Bucket:    os.Getenv("MINIO_BUCKET"),
			Region:    os.Getenv("MINIO_REGION"),
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			SecretKey: os.Getenv("MINIO_SECRET_KEY"),
			Endpoint:  os.Getenv("MINIO_ENDPOINT"),
			Prefixes: BucketPrefix{
				Mamography:        getEnvOrDefault("MINIO_PREFIX_MAMOGRAPHY", "mamography"),
				Cancer:            getEnvOrDefault("MINIO_PREFIX_CANCER", "cancer"),
				GeneticTest:       getEnvOrDefault("MINIO_PREFIX_GENETIC_TEST", "genetic-test"),
				FatherGeneticTest: getEnvOrDefault("MINIO_PREFIX_FATHER_GENETIC_TEST", "father-genetic-test"),
				MotherGeneticTest: getEnvOrDefault("MINIO_PREFIX_MOTHER_GENETIC_TEST", "mother-genetic-test"),
			},
		},
		CalcURL: CalcURL{
			Premm5: os.Getenv("PREMM5_API_URL"),
			BCRA:   os.Getenv("BCRA_API_URL"),
			Gail:   os.Getenv("GAIL_API_URL"),
			PLCO:   os.Getenv("PLCO_API_URL"),
		},
		VerificationAPI: VerificationAPI{
			BaseURL: os.Getenv("VERIFICATION_API_URL"),
			APIKey:  os.Getenv("VERIFICATION_API_KEY"),
		},
		Security: Security{
			EncryptionKey:      os.Getenv("ENCRYPTION_KEY"),
			RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 5),
			RateLimitWindow:    getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 1),
		},
		JWT: JWT{
			PrivateKey: getEnvPEM("JWT_PRIVATE_KEY"),
			PublicKey:  getEnvPEM("JWT_PUBLIC_KEY"),
		},
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvPEM(key string) string {
	return strings.ReplaceAll(os.Getenv(key), `\n`, "\n")
}
