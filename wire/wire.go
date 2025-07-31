//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/middleware"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDatabase,
	wire.Bind(new(database.Database), new(*database.PostgresDatabase)),
	wire.Struct(new(Database), "*"),
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

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	MiddlewareProviderSet,
	ProvideDBConfig,
	ProvideConstants,
)

type Database struct {
	DB database.Database
}

type Middlewares struct {
	Recovery *middleware.RecoveryMiddleware
}

type Application struct {
	Database    *Database
	Middlewares *Middlewares
}

func NewApplication(
	database *Database,
	middlewares *Middlewares,
) *Application {
	return &Application{
		Database:    database,
		Middlewares: middlewares,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
