package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kodia-studio/kodia/internal/infrastructure/cache"
	"github.com/kodia-studio/kodia/internal/infrastructure/database"
	"github.com/kodia-studio/kodia/internal/infrastructure/logger"
	"github.com/kodia-studio/kodia/internal/infrastructure/worker"
	events_infra "github.com/kodia-studio/kodia/internal/infrastructure/events"
	"github.com/kodia-studio/kodia/internal/core/events"
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

	log.Info("Starting Kodia Worker",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	// 3. Initialize database (for use in jobs)
	_, err = database.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Initialize cache (Redis)
	cacheProvider, err := cache.New(cfg, log)
	if err != nil {
		log.Fatal("Redis connection failed, worker requires Redis", zap.Error(err))
	}
	_ = cacheProvider // Avoid unused var if not injected yet

	// 5. Initialize Worker Processor
	processor := worker.NewProcessor(cfg, log)

	// 5.1 Initialize Dispatcher for Worker (to find listener instances)
	queueProvider := worker.NewAsynqProvider(cfg, log)
	dispatcher := events_infra.NewDispatcher(queueProvider, log)
	events.RegisterEvents(dispatcher)

	// Register generic event listener handler
	processor.Register("event.listener.job", dispatcher.HandleListenerTask)

	// --- Job Registration Start ---
	// Jobs will be automatically registered here by the CLI
	// --- Job Registration End ---

	// 6. Start Worker with premature shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := processor.Start(); err != nil {
			log.Fatal("Worker error", zap.Error(err))
		}
	}()

	log.Info("Worker is running and listening for jobs...")

	<-quit
	log.Info("Shutting down worker...")
}
