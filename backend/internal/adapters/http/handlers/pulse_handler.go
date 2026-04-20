package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kodia-studio/kodia/pkg/observability"
	"go.uber.org/zap"
)

// PulseHandler handles real-time monitoring WebSocket connections.
type PulseHandler struct {
	manager  *observability.PulseManager
	log      *zap.Logger
	upgrader websocket.Upgrader
}

// NewPulseHandler creates a new PulseHandler instance.
func NewPulseHandler(manager *observability.PulseManager, log *zap.Logger) *PulseHandler {
	return &PulseHandler{
		manager: manager,
		log:     log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development, refine for production
				return true
			},
		},
	}
}

// Stream handles the WebSocket connection for Pulse.
// Protected by Admin middleware.
func (h *PulseHandler) Stream(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("Failed to upgrade Pulse WebSocket", zap.Error(err))
		return
	}
	defer conn.Close()

	h.log.Info("Pulse WebSocket client connected")

	// Register client with manager
	clientChan := h.manager.Register()
	defer h.manager.Unregister(clientChan)

	// Ping-pong and read loop to keep connection alive and detect disconnects
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	// Write loop
	for message := range clientChan {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.log.Warn("Failed to send Pulse message", zap.Error(err))
			break
		}
	}

	h.log.Info("Pulse WebSocket client disconnected")
}
