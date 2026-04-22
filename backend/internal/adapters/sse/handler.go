package sse

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles SSE HTTP connections.
type Handler struct {
	manager *Manager
	log     *zap.Logger
}

// NewHandler creates a new SSE Handler.
func NewHandler(manager *Manager, log *zap.Logger) *Handler {
	return &Handler{manager: manager, log: log}
}

// ServePublic handles a public SSE stream for a given channel.
// GET /api/v1/sse/:channel
func (h *Handler) ServePublic(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		channel = "public"
	}
	h.stream(c, "", channel)
}

// ServeUser handles a private SSE stream for the authenticated user.
// GET /api/v1/sse/user
// Requires user_id in gin context (set by Auth middleware).
func (h *Handler) ServeUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(string)
	h.stream(c, uid, fmt.Sprintf("private-%s", uid))
}

func (h *Handler) stream(c *gin.Context, userID, channel string) {
	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	client := h.manager.Subscribe(userID, channel)
	defer h.manager.Unsubscribe(client.ID)

	h.log.Info("SSE client connected",
		zap.String("client_id", client.ID),
		zap.String("user_id", userID),
		zap.String("channel", channel),
	)

	// Send initial connection event
	connected := &SSEEvent{
		Event: "connected",
		Data:  map[string]string{"channel": channel, "client_id": client.ID},
	}
	fmt.Fprint(c.Writer, connected.Format())
	c.Writer.Flush()

	// Heartbeat ticker to keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	notify := c.Request.Context().Done()
	w := c.Writer

	c.Stream(func(w io.Writer) bool {
		select {
		case <-notify:
			// Client disconnected
			h.log.Info("SSE client disconnected", zap.String("client_id", client.ID))
			return false

		case event, ok := <-client.Events:
			if !ok {
				return false
			}
			fmt.Fprint(w, event.Format())
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			return true

		case <-ticker.C:
			// Send heartbeat comment
			fmt.Fprintf(w, ": heartbeat\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			return true
		}
	})
	_ = w
}

// Status returns current SSE stats.
// GET /api/v1/sse/status
func (h *Handler) Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"active_connections": h.manager.ClientCount(),
	})
}
