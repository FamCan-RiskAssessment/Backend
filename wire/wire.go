//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDatabase,
	wire.Bind(new(database.Database), new(*database.PostgresDatabase)),
	wire.Struct(new(Database), "*"),
)

func ProvideDBConfig(container *bootstrap.Config) *bootstrap.Database {
	return &container.Env.Database
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	ProvideDBConfig,
)

type Database struct {
	DB database.Database
}

type Application struct {
	Database *Database
}

func NewApplication(
	database *Database,
) *Application {
	return &Application{
		Database: database,
	}
}

func InitializeApplication(config *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
