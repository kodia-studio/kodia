// Package http provides the router setup for Kodia Framework.
package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Router holds all dependencies for the HTTP router.
type Router struct {
	cfg         *config.Config
	log         *zap.Logger
	jwtManager  *jwt.Manager
	authHandler *handlers.AuthHandler
	userHandler *handlers.UserHandler
	redisClient *redis.Client
}

// NewRouter creates a new Router instance.
func NewRouter(
	cfg *config.Config,
	log *zap.Logger,
	jwtManager *jwt.Manager,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	redisClient *redis.Client,
) *Router {
	return &Router{
		cfg:         cfg,
		log:         log,
		jwtManager:  jwtManager,
		authHandler: authHandler,
		userHandler: userHandler,
		redisClient: redisClient,
	}
}

// Setup configures the Gin engine with middleware and routes.
func (r *Router) Setup() *gin.Engine {
	if r.cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Global middleware
	engine.Use(middleware.Recovery(r.log))
	engine.Use(middleware.Logger(r.log))

	// Validate CORS configuration for security issues
	if err := ValidateCORSConfig(r.cfg, r.log); err != nil {
		r.log.Fatal("CORS configuration validation failed", zap.Error(err))
	}

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     r.cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	engine.Use(cors.New(corsConfig))

	// API grouped routes
	api := engine.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Kodia Backend is healthy",
			})
		})

		// Auth routes with rate limiting
		auth := api.Group("/auth")
		{
			// Apply rate limiting middleware to auth endpoints if Redis is available
			if r.redisClient != nil {
				authLimiter := middleware.AuthEndpointRateLimiter(r.redisClient, r.log)
				auth.POST("/register", authLimiter.Middleware(), r.authHandler.Register)
				auth.POST("/login", authLimiter.Middleware(), r.authHandler.Login)
				auth.POST("/refresh", authLimiter.Middleware(), r.authHandler.RefreshToken)
			} else {
				// No rate limiting if Redis is not available
				r.log.Warn("Redis client not available, rate limiting disabled for auth endpoints")
				auth.POST("/register", r.authHandler.Register)
				auth.POST("/login", r.authHandler.Login)
				auth.POST("/refresh", r.authHandler.RefreshToken)
			}
			auth.POST("/logout", r.jwtManagerAuthMiddleware(), r.authHandler.Logout)

			// Protected auth routes
			protectedAuth := auth.Group("")
			protectedAuth.Use(r.jwtManagerAuthMiddleware())
			{
				protectedAuth.POST("/logout-all", r.authHandler.LogoutAll)
				protectedAuth.GET("/me", r.authHandler.Me)
			}
		}

		// User routes
		users := api.Group("/users")
		users.Use(r.jwtManagerAuthMiddleware())
		{
			users.GET("/me", r.userHandler.GetMe)
			users.POST("/me/change-password", r.userHandler.ChangePassword)

			// Admin only routes
			adminUsers := users.Group("")
			adminUsers.Use(middleware.RequireRole("admin"))
			{
				adminUsers.GET("", r.userHandler.GetAll)
				adminUsers.GET("/:id", r.userHandler.GetByID)
				adminUsers.PATCH("/:id", r.userHandler.Update)
				adminUsers.DELETE("/:id", r.userHandler.Delete)
			}
		}
	}

	return engine
}

func (r *Router) jwtManagerAuthMiddleware() gin.HandlerFunc {
	return middleware.Auth(r.jwtManager)
}

func formatValidationErrors(err error) map[string][]string {
	var ve validator.ValidationErrors
	errs := make(map[string][]string)
	for _, fe := range ve {
		errs[fe.Field()] = append(errs[fe.Field()], fe.Tag())
	}
	return errs
}
