package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kodia/framework/backend/internal/adapters/http/handlers"
	kodia_http "github.com/kodia/framework/backend/internal/adapters/http"
	"github.com/kodia/framework/backend/internal/adapters/repository/postgres"
	"github.com/kodia/framework/backend/internal/core/services"
	"github.com/kodia/framework/backend/internal/infrastructure/cache"
	"github.com/kodia/framework/backend/internal/infrastructure/database"
	"github.com/kodia/framework/backend/internal/infrastructure/logger"
	"github.com/kodia/framework/backend/pkg/config"
	"github.com/kodia/framework/backend/pkg/jwt"
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
	_, err = cache.New(cfg, log)
	if err != nil {
		log.Warn("Cache (Redis) connection failed, some features may be limited", zap.Error(err))
	}

	// 5. Initialize JWT manager
	jwtManager := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiryHours,
		cfg.JWT.RefreshExpiryDays,
	)

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
