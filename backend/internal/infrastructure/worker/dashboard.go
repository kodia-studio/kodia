package worker

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/kodia-studio/kodia/pkg/config"
)

// DashboardHandler exposes the Asynqmon dashboard.
type DashboardHandler struct {
	monHandler *asynqmon.HTTPHandler
}

// NewDashboardHandler creates a new monitoring dashboard handler.
func NewDashboardHandler(cfg *config.Config) *DashboardHandler {
	h := asynqmon.New(asynqmon.Options{
		RootPath: "/api/admin/queues", // Must match router mount point
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
	})

	return &DashboardHandler{monHandler: h}
}

// Mount mounts the dashboard to a Gin engine.
func (h *DashboardHandler) Mount(r *gin.Engine) {
	// asynqmon uses gorilla/mux internally, we wrap it for Gin
	r.Any("/api/admin/queues/*any", func(c *gin.Context) {
		h.monHandler.ServeHTTP(c.Writer, c.Request)
	})
}

// FailedJobsAPI (Future): Custom endpoint for PostgreSQL auditing if needed.
func (h *DashboardHandler) GetFailedJobs(c *gin.Context) {
	// Implementation for PostgreSQL fallback audit logs
	c.JSON(http.StatusOK, gin.H{"message": "Audit logs feature coming soon"})
}
