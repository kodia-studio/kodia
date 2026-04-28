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

// Live godoc
// @Summary      Liveness probe
// @Description  Returns 200 if service is alive (simple ping)
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

// Ready godoc
// @Summary      Readiness probe
// @Description  Checks service readiness including database and Redis connectivity
// @Tags         health
// @Produce      json
// @Success      200 {object} health.Stats "Service is ready"
// @Failure      503 {object} health.Stats "Service is degraded or down"
// @Router       /health/ready [get]
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
