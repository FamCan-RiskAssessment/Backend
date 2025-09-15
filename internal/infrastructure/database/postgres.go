package database

import (
	"fmt"
	"sync"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
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

func NewPostgresDatabase(dbConfig *bootstrap.Database) *PostgresDatabase {
	fmt.Println(dbConfig.Host)
	fmt.Println(dbConfig.Port)
	fmt.Println(dbConfig.User)
	fmt.Println(dbConfig.Password)
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
