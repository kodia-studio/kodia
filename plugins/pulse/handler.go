package pulse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Handler struct {
	manager  *Manager
	log      *zap.Logger
	upgrader websocket.Upgrader
}

func NewHandler(manager *Manager, log *zap.Logger) *Handler {
	return &Handler{
		manager: manager,
		log:     log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handler) Stream(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("Failed to upgrade Pulse WebSocket", zap.Error(err))
		return
	}
	defer conn.Close()

	h.log.Info("Pulse WebSocket client connected")

	clientChan := h.manager.Register()
	defer h.manager.Unregister(clientChan)

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	for message := range clientChan {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.log.Warn("Failed to send Pulse message", zap.Error(err))
			break
		}
	}

	h.log.Info("Pulse WebSocket client disconnected")
}
