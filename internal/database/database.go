package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/eduardocardona93/golang_shopping/internal/config"
	"github.com/eduardocardona93/golang_shopping/internal/domain"
)

// Connect opens a GORM connection to PostgreSQL, verifies it with a ping and
// configures the underlying connection pool.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := gormlogger.Warn
	if cfg.AppEnv == "development" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB handle: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return db, nil
}

// AutoMigrate creates/updates the database schema for all domain models.
func AutoMigrate(db *gorm.DB) error {
	log.Println("running database auto-migration...")
	return db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Inventory{},
		&domain.Purchase{},
		&domain.PurchaseItem{},
	)
}
