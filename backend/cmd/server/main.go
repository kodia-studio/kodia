package main

import (
	"fmt"
	"os"

	"github.com/kodia-studio/kodia/internal/infrastructure/logger"
	"github.com/kodia-studio/kodia/internal/providers"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/kodia-studio/kodia/pkg/kodia"
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

	log.Info("Starting Kodia Framework",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	// 3. Initialize Kodia App
	app := kodia.NewApp(cfg, log)

	// 4. Register Official Providers (Batteries Included)
	err = app.RegisterProviders(
		providers.NewDatabaseProvider(),
		providers.NewInfraProvider(),
		providers.NewHttpProvider(),
		providers.NewAuthProvider(),
		providers.NewUserProvider(),
		providers.NewNotificationProvider(),
		providers.NewWebSocketProvider(),
		providers.NewGraphQLProvider(),
		providers.NewPulseProvider(), // Broadcaster for real-time monitoring
		// Third-party plugins would be added here
	)
	if err != nil {
		log.Fatal("Failed to register providers", zap.Error(err))
	}

	// 5. Boot all providers
	if err := app.Boot(); err != nil {
		log.Fatal("Failed to boot providers", zap.Error(err))
	}

	// 6. Run the application
	if err := app.Run(); err != nil {
		log.Fatal("Application execution failed", zap.Error(err))
	}
}
