package database

import (
	"fmt"
	"sync"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
	WithTransaction(fn func(Database) error) error
}

type PostgresDatabase struct {
	DB *gorm.DB
}

var (
	dbOnce     sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(dbConfig *bootstrap.Database, server *bootstrap.Server) *PostgresDatabase {
	dbOnce.Do(func() {
		// TODO: Security Enhancement - Enable SSL/TLS for database connections
		// Currently using sslmode=disable which is insecure for production
		// Action needed:
		// 1. Configure PostgreSQL server with SSL certificate
		// 2. Change sslmode=disable to sslmode=require
		// 3. Optionally add sslcert, sslkey, and sslrootcert parameters for mutual TLS
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logging.GormLogger(server.Mode),
		})
		if err != nil {
			panic(fmt.Errorf("failed to connect to database"))
		}

		dbInstance = &PostgresDatabase{DB: db}
	})
	return dbInstance
}

func (db *PostgresDatabase) GetDB() *gorm.DB {
	return dbInstance.DB
}

func (pgx *PostgresDatabase) WithTransaction(fn func(Database) error) error {
	return pgx.DB.Transaction(func(tx *gorm.DB) error {
		txWrapper := &PostgresDatabase{DB: tx}
		return fn(txWrapper)
	})
}
