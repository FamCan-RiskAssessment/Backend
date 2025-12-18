//go:build wireinject
// +build wireinject

package wire

import (
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/service"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/communication"
	domainExternal "github.com/FamCan-RiskAssessment/Backend/internal/domain/external"
	domainPostgre "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	domainRedis "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/redis"
	domainS3 "github.com/FamCan-RiskAssessment/Backend/internal/domain/storage/s3"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/communication/sms"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/crypto"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	infraExternal "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/external"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/ratelimit"
	infraJWT "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/jwt"
	infraLocalization "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/localization"
	infraPostgre "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/repository/postgres"
	infraRedis "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/repository/redis"
	seed "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/seed"
	infraStorage "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/storage"
	actionlog "github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller/action_log"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller/calc"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller/form"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller/user"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/middleware"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDatabase,
	database.NewRedisDatabase,
	wire.Bind(new(database.Database), new(*database.PostgresDatabase)),
	wire.Bind(new(database.Cache), new(*database.RedisDatabase)),
	wire.Struct(new(Database), "*"),
)

var RepositoryProviderSet = wire.NewSet(
	infraPostgre.NewUserRepository,
	infraPostgre.NewFormRepository,
	infraPostgre.NewActionLogRepository,
	infraRedis.NewUserCacheRepository,
	wire.Bind(new(domainPostgre.UserRepository), new(*infraPostgre.UserRepository)),
	wire.Bind(new(domainPostgre.FormRepository), new(*infraPostgre.FormRepository)),
	wire.Bind(new(domainPostgre.ActionLogRepository), new(*infraPostgre.ActionLogRepository)),
	wire.Bind(new(domainRedis.UserCacheRepository), new(*infraRedis.UserCacheRepository)),
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewFormService,
	service.NewJWTService,
	service.NewOTPService,
	service.NewActionLogService,
	service.NewCalcService,
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.FormService), new(*service.FormService)),
	wire.Bind(new(usecase.OtpService), new(*service.OTPService)),
	wire.Bind(new(usecase.JwtService), new(*service.JWTService)),
	wire.Bind(new(usecase.ActionLogService), new(*service.ActionLogService)),
	wire.Bind(new(usecase.CalcService), new(*service.CalcService)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	form.NewGeneralFormController,
	wire.Struct(new(GeneralControllers), "*"),
)

var AdminControllerProviderSet = wire.NewSet(
	user.NewAdminUserController,
	form.NewAdminFormController,
	actionlog.NewActionLogController,
	calc.NewAdminCalcController,
	wire.Struct(new(AdminControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	form.NewCustomerFormController,
	wire.Struct(new(CustomerControllers), "*"),
)

var ControllerProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var AdapterProviderSet = wire.NewSet(
	infraJWT.NewJWTKeyManager,
	infraLocalization.NewTranslationService,
	infraStorage.NewS3Storage,
	sms.NewAsanakSMSService,
	infraExternal.NewVerificationClient,
	wire.Bind(new(domainS3.S3Storage), new(*infraStorage.S3Storage)),
	wire.Bind(new(communication.SmsService), new(*sms.AsanakSMSService)),
	wire.Bind(new(domainExternal.VerificationClient), new(*infraExternal.VerificationClientImpl)),
)

var RateLimitProviderSet = wire.NewSet(
	middleware.NewRateLimitMiddleware,
	ProvideRateLimiter,
)

var CryptoProviderSet = wire.NewSet(
	crypto.NewPasswordHasher,
	crypto.NewFieldEncryptor,
	wire.Bind(new(usecase.PasswordHasher), new(*crypto.PasswordHasher)),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewCorsMiddleware,
	middleware.NewRecoveryMiddleware,
	middleware.NewLocalizationMiddleware,
	middleware.NewAuthMiddleware,
	RateLimitProviderSet,
	wire.Struct(new(Middlewares), "*"),
)

var SeedProviderSet = wire.NewSet(
	seed.NewRoleSeeder,
	seed.NewDummySeeder,
	wire.Struct(new(Seeds), "*"),
)

func ProvideDBConfig(container *bootstrap.Config) *bootstrap.Database {
	return &container.Env.Database
}

func ProvideConstants(container *bootstrap.Config) *bootstrap.Constants {
	return container.Constants
}

func ProvideRDBConfig(container *bootstrap.Config) *bootstrap.Redis {
	return &container.Env.Cache
}

func ProvideOTPConfig(container *bootstrap.Config) *bootstrap.OTP {
	return &container.Env.OTP
}

func ProvideStorageConfig(container *bootstrap.Config) *bootstrap.S3 {
	return &container.Env.S3
}

func ProvideSMSGatewayConfig(container *bootstrap.Config) *bootstrap.SMSGateway {
	return &container.Env.SMSGateway
}

func ProvideSMSTemplates(container *bootstrap.Config) *bootstrap.SMSTemplates {
	return &container.Constants.SMSTemplates
}

func ProvideJWTKeysPath(container *bootstrap.Config) *bootstrap.JWTKeysPath {
	return &container.Constants.JWTKeysPath
}

func ProvideSuperAdminCredentials(container *bootstrap.Config) *bootstrap.SuperAdmin {
	return &container.Env.SuperAdmin
}

func ProvidePagination(container *bootstrap.Config) *bootstrap.Pagination {
	return &container.Env.Pagination
}

func ProvideCalcURL(container *bootstrap.Config) *bootstrap.CalcURL {
	return &container.Env.CalcURL
}

func ProvideVerificationAPIConfig(container *bootstrap.Config) *bootstrap.VerificationAPI {
	return &container.Env.VerificationAPI
}

func ProvideSecurityConfig(container *bootstrap.Config) *bootstrap.Security {
	return &container.Env.Security
}

func ProvideRateLimiter(rdb database.Cache, security *bootstrap.Security) *ratelimit.RateLimiter {
	return ratelimit.NewRateLimiter(
		rdb.GetRDB(),
		security.RateLimitPerMinute,
		time.Duration(security.RateLimitWindow)*time.Second,
	)
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	MiddlewareProviderSet,
	GeneralControllerProviderSet,
	AdminControllerProviderSet,
	CustomerControllerProviderSet,
	ControllerProviderSet,
	AdapterProviderSet,
	CryptoProviderSet,
	ProvideDBConfig,
	ProvideConstants,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideStorageConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideJWTKeysPath,
	ProvideSuperAdminCredentials,
	ProvidePagination,
	ProvideCalcURL,
	ProvideVerificationAPIConfig,
	ProvideSecurityConfig,
	SeedProviderSet,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController *user.GeneralUserController
	FormController *form.GeneralFormController
}

type AdminControllers struct {
	UserController      *user.AdminUserController
	FormController      *form.AdminFormController
	ActionLogController *actionlog.ActionLogController
	CalcController      *calc.AdminCalcController
}

type CustomerControllers struct {
	FormController *form.CustomerFormController
}

type Controllers struct {
	General  *GeneralControllers
	Admin    *AdminControllers
	Customer *CustomerControllers
}

type Middlewares struct {
	Cors         *middleware.CORSMiddleware
	Recovery     *middleware.RecoveryMiddleware
	Localization *middleware.LocalizationMiddleware
	Auth         *middleware.AuthMiddleware
	RateLimit    *middleware.RateLimitMiddleware
}

type Seeds struct {
	RoleSeeder  *seed.RoleSeeder
	DummySeeder *seed.DummySeeder
}

type Application struct {
	Database    *Database
	Middlewares *Middlewares
	Controllers *Controllers
	Seeds       *Seeds
}

func NewApplication(
	database *Database,
	middlewares *Middlewares,
	controllers *Controllers,
	seeds *Seeds,
) *Application {
	return &Application{
		Database:    database,
		Middlewares: middlewares,
		Controllers: controllers,
		Seeds:       seeds,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
