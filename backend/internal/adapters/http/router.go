// Package http provides the router setup for Kodia Framework.
package http

import (
	"net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/adapters/http/handlers"
	"github.com/kodia-studio/kodia/internal/adapters/http/middleware"
	"github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"github.com/kodia-studio/kodia/pkg/observability"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/kodia-studio/kodia/docs" // Swagger docs
	"go.uber.org/zap"
)

// Router holds all dependencies for the HTTP router.
type Router struct {
	cfg         *config.Config
	log         *zap.Logger
	jwtManager  *jwt.Manager
	authHandler *handlers.AuthHandler
	userHandler *handlers.UserHandler
	redisClient    *redis.Client
	wsHandler      *websocket.Handler
	graphqlHandler *handlers.GraphQLHandler
	obsManager     *observability.Manager
	pulseHandler   *handlers.PulseHandler
	healthHandler  *handlers.HealthHandler
}

// NewRouter creates a new Router instance.
func NewRouter(
	cfg *config.Config,
	log *zap.Logger,
	jwtManager *jwt.Manager,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	redisClient *redis.Client,
	wsHandler *websocket.Handler,
	graphqlHandler *handlers.GraphQLHandler,
	obsManager *observability.Manager,
	pulseHandler *handlers.PulseHandler,
	healthHandler *handlers.HealthHandler,
) *Router {
	return &Router{
		cfg:            cfg,
		log:            log,
		jwtManager:     jwtManager,
		authHandler:    authHandler,
		userHandler:    userHandler,
		redisClient:    redisClient,
		wsHandler:      wsHandler,
		graphqlHandler: graphqlHandler,
		obsManager:     obsManager,
		pulseHandler:   pulseHandler,
		healthHandler:  healthHandler,
	}
}

// Setup configures the Gin engine with middleware and routes.
func (r *Router) Setup() *gin.Engine {
	if r.cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Global middleware
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery(r.log))
	engine.Use(middleware.Logger(r.log))

	// Tracing & Metrics Middleware
	if r.cfg.Observability.TracingEnabled {
		engine.Use(middleware.Tracing(r.cfg.Observability.ServiceName))
	}
	if r.cfg.Observability.MetricsEnabled {
		engine.Use(middleware.Metrics())
	}

	// Validate CORS configuration for security issues
	if err := ValidateCORSConfig(r.cfg, r.log); err != nil {
		r.log.Fatal("CORS configuration validation failed", zap.Error(err))
	}

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     r.cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Trace-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Trace-ID"},
		AllowCredentials: true,
	}
	engine.Use(cors.New(corsConfig))

	// Profiling endpoints (pprof) - Development only
	if !r.cfg.IsProduction() {
		pprofGroup := engine.Group("/debug/pprof")
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
			pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
			pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
			pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
			pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
			pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
		}
	}

	// API grouped routes
	api := engine.Group("/api")
	{
		// Apply global rate limiting to all /api/* routes
		if r.redisClient != nil {
			globalLimiter := middleware.LooseRateLimiter(r.redisClient, r.log)
			api.Use(globalLimiter.Middleware())
		}

		// API Documentation (Swagger) - Development only
		if !r.cfg.IsProduction() {
			api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		v1 := api.Group("/v1")
		{
			// Health check (Enhanced)
			v1.GET("/health", r.healthHandler.Ready)
			v1.GET("/health/live", r.healthHandler.Live)
			v1.GET("/health/ready", r.healthHandler.Ready)

			// Auth routes with rate limiting
			auth := v1.Group("/auth")
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
			users := v1.Group("/users")
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

			// WebSocket routes
			ws := v1.Group("ws")
			{
				ws.GET("", r.wsHandler.ServeWS)
				ws.GET("/room/:room", r.wsHandler.ServeRoom)
				ws.GET("/status", r.wsHandler.GetStatus)
			}

			// Pulse (Monitoring) routes
			pulse := v1.Group("/pulse")
			pulse.Use(r.jwtManagerAuthMiddleware())
			pulse.Use(middleware.RequireRole("admin"))
			{
				pulse.GET("/stream", r.pulseHandler.Stream)
			}

			// GraphQL routes
			v1.POST("/query", middleware.GraphQLContextMiddleware(), r.graphqlHandler.QueryHandler())
			if !r.cfg.IsProduction() {
				v1.GET("/playground", r.graphqlHandler.PlaygroundHandler())
			}
		}
	}

	return engine
}

func (r *Router) jwtManagerAuthMiddleware() gin.HandlerFunc {
	return middleware.Auth(r.jwtManager)
}

