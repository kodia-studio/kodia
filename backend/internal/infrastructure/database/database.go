// Package database provides GORM database connection for Kodia Framework.
// Supports PostgreSQL and MySQL via the configured driver.
package database

import (
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New establishes and returns a GORM database connection.
// The driver is determined by cfg.Database.Driver ("postgres" or "mysql").
func New(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Database.DSN())
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.Database.DSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (supported: postgres, mysql)", cfg.Database.Driver)
	}

	logLevel := gormlogger.Silent
	if cfg.IsDevelopment() {
		logLevel = gormlogger.Info
	}

	log.Debug("Attempting database connection",
		zap.String("driver", cfg.Database.Driver),
		zap.String("database", cfg.Database.Name),
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Connection pool settings
	maxOpen := cfg.Database.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 25
	}
	maxIdle := cfg.Database.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 10
	}
	lifetime := cfg.Database.ConnMaxLifetime
	if lifetime == 0 {
		lifetime = 30 * time.Minute
	}

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(lifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info("Database connected",
		zap.String("driver", cfg.Database.Driver),
		zap.String("database", cfg.Database.Name),
	)

	return db, nil
}
