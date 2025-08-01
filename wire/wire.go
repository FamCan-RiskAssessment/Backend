//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/service"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/communication"
	domainPostgre "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	domainRedis "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/redis"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/communication/sms"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	infraJWT "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/jwt"
	infraPostgre "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/repository/postgres"
	infraRedis "github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/repository/redis"
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
	infraRedis.NewUserCacheRepository,
	wire.Bind(new(domainPostgre.UserRepository), new(*infraPostgre.UserRepository)),
	wire.Bind(new(domainRedis.UserCacheRepository), new(*infraRedis.UserCacheRepository)),
)

var ServiceProviderSet = wire.NewSet(
	service.NewUserService,
	service.NewJWTService,
	service.NewOTPService,
	sms.NewSMSService,
	wire.Bind(new(usecase.UserService), new(*service.UserService)),
	wire.Bind(new(usecase.OtpService), new(*service.OTPService)),
	wire.Bind(new(usecase.JwtService), new(*service.JWTService)),
	wire.Bind(new(communication.SmsService), new(*sms.SMSService)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	wire.Struct(new(GeneralControllers), "*"),
)

var ControllerProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var AdapterProviderSet = wire.NewSet(
	infraJWT.NewJWTKeyManager,
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewRecoveryMiddleware,
	wire.Struct(new(Middlewares), "*"),
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

func ProvideSMSGatewayConfig(container *bootstrap.Config) *bootstrap.SMSGateway {
	return &container.Env.SMSGateway
}

func ProvideSMSTemplates(container *bootstrap.Config) *bootstrap.SMSTemplates {
	return &container.Constants.SMSTemplates
}

func ProvideJWTKeysPath(container *bootstrap.Config) *bootstrap.JWTKeysPath {
	return &container.Constants.JWTKeysPath
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	MiddlewareProviderSet,
	GeneralControllerProviderSet,
	ControllerProviderSet,
	AdapterProviderSet,
	ProvideDBConfig,
	ProvideConstants,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideJWTKeysPath,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController *user.GeneralUserController
}

type Controllers struct {
	General *GeneralControllers
}

type Middlewares struct {
	Recovery *middleware.RecoveryMiddleware
}

type Application struct {
	Database    *Database
	Middlewares *Middlewares
	Controllers *Controllers
}

func NewApplication(
	database *Database,
	middlewares *Middlewares,
	controllers *Controllers,
) *Application {
	return &Application{
		Database:    database,
		Middlewares: middlewares,
		Controllers: controllers,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
