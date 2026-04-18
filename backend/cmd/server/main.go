package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	kodia_http "github.com/kodia-studio/kodia/internal/adapters/http"
	"github.com/kodia-studio/kodia/internal/adapters/repository/postgres"
	"github.com/kodia-studio/kodia/internal/core/services"
	"github.com/kodia-studio/kodia/internal/infrastructure/cache"
	"github.com/kodia-studio/kodia/internal/infrastructure/database"
	"github.com/kodia-studio/kodia/internal/infrastructure/logger"
	"github.com/kodia-studio/kodia/internal/infrastructure/storage"
	"github.com/kodia-studio/kodia/internal/infrastructure/worker"
	"github.com/kodia-studio/kodia/internal/infrastructure/mailer"
	events_infra "github.com/kodia-studio/kodia/internal/infrastructure/events"
	"github.com/kodia-studio/kodia/internal/core/events"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/kodia-studio/kodia/pkg/jwt"
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

	log.Info("Starting Kodia Backend",
		zap.String("app_name", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port),
	)

	// 3. Initialize database
	db, err := database.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Initialize cache (Redis)
	cacheProvider, err := cache.New(cfg, log)
	if err != nil {
		log.Warn("Cache (Redis) connection failed, some features may be limited", zap.Error(err))
	}
	_ = cacheProvider // Avoid unused var if not injected yet

	// 5. Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiryHours,
		cfg.JWT.RefreshExpiryDays,
	)

	// 5.1 Initialize Storage Provider
	var storageProvider ports.StorageProvider
	switch strings.ToLower(cfg.Storage.Provider) {
	case "s3":
		var err error
		storageProvider, err = storage.NewS3StorageProvider(cfg)
		if err != nil {
			log.Fatal("Failed to initialize S3 storage", zap.Error(err))
		}
		log.Info("S3 Storage initialized", zap.String("bucket", cfg.Storage.Bucket))
	default:
		storageProvider = storage.NewLocalStorageProvider(cfg)
		log.Info("Local Storage initialized", zap.String("dir", cfg.Storage.LocalDir))
	}
	_ = storageProvider // Avoid unused var if not injected yet

	// 5.2 Initialize Mailer
	mailProvider, err := mailer.NewSMTPMailer(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize Mailer", zap.Error(err))
	}
	log.Info("SMTP Mailer initialized", zap.String("host", cfg.Mail.Host))
	_ = mailProvider // Avoid unused var if not injected yet

	// 5.3 Initialize Event Dispatcher
	queueProvider := worker.NewAsynqProvider(cfg, log)
	dispatcher := events_infra.NewDispatcher(queueProvider, log)
	events.RegisterEvents(dispatcher)
	log.Info("Event Dispatcher initialized and events registered")
	_ = dispatcher // Avoid unused var if not injected yet

	// 6. Initialize validation
	validate := validator.New()

	// 7. Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)

	// 8. Initialize services
	authService := services.NewAuthService(userRepo, refreshRepo, jwtManager, log)
	userService := services.NewUserService(userRepo, log)

	// 9. Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, validate, log)
	userHandler := handlers.NewUserHandler(userService, validate, log)

	// 10. Initialize router
	router := kodia_http.NewRouter(cfg, log, jwtManager, authHandler, userHandler)
	engine := router.Setup()

	// 11. Start server with graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Listen and serve error", zap.Error(err))
		}
	}()

	log.Info(fmt.Sprintf("Server is running on http://localhost:%d", cfg.App.Port))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need added it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exiting gracefully")
}
