package main

import (
	"fmt"
	"os"

	"github.com/kodia-studio/kodia/internal/infrastructure/database"
	"github.com/kodia-studio/kodia/internal/infrastructure/database/seeders"
	"github.com/kodia-studio/kodia/internal/infrastructure/logger"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize logger
	log, err := logger.New(cfg.IsDevelopment())
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Database Seeder",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	// 3. Initialize database
	db, err := database.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Run all seeders
	log.Info("Running all registered seeders...")
	if err := seeders.RunAll(db); err != nil {
		log.Fatal("Seeding failed", zap.Error(err))
	}

	log.Info("Database seeding completed successfully! ✅")
}
