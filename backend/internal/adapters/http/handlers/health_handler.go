package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/health"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// HealthHandler handles system health check requests.
type HealthHandler struct {
	db    *gorm.DB
	redis *redis.Client
	log   *zap.Logger
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(db *gorm.DB, redis *redis.Client, log *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
		log:   log,
	}
}

// Live handles liveness check requests (/health/live).
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

// Ready handles readiness check requests (/health/ready).
func (h *HealthHandler) Ready(c *gin.Context) {
	checkers := []health.Checker{
		&health.DBChecker{DB: h.db},
	}

	if h.redis != nil {
		checkers = append(checkers, &health.RedisChecker{Client: h.redis})
	}

	stats, _ := health.Gather(c.Request.Context(), checkers...)

	status := http.StatusOK
	if stats.Status == "degraded" || stats.Status == "down" {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"success": stats.Status == "up",
		"data":    stats,
	})
}
