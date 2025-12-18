package bootstrap

import (
	"os"
	"strconv"

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
	S3              S3
	CalcURL         CalcURL
	VerificationAPI VerificationAPI
	Security        Security
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

type S3 struct {
	Buckets   BucketName
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type BucketName struct {
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
		S3: S3{
			Region:    os.Getenv("S3_REGION"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Buckets: BucketName{
				Mamography:        os.Getenv("MAMOGRAPHY_BUCKETNAME"),
				Cancer:            os.Getenv("CANCER_BUCKETNAME"),
				GeneticTest:       os.Getenv("GENETIC_TEST_BUCKETNAME"),
				FatherGeneticTest: os.Getenv("FATHER_GENETIC_TEST_BUCKETNAME"),
				MotherGeneticTest: os.Getenv("MOTHER_GENETIC_TEST_BUCKETNAME"),
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
