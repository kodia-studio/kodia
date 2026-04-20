package websocket

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kodia-studio/kodia/pkg/jwt"
	"go.uber.org/zap"
)

// Handler handles WebSocket connections and upgrades from HTTP.
type Handler struct {
	hub        *Hub
	jwtManager *jwt.Manager
	log        *zap.Logger
	upgrader   websocket.Upgrader
}

// NewHandler creates a new WebSocket Handler.
func NewHandler(hub *Hub, jwtManager *jwt.Manager, log *zap.Logger) *Handler {
	return &Handler{
		hub:        hub,
		jwtManager: jwtManager,
		log:        log,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// ServeWS handles WebSocket upgrades at the /api/ws endpoint.
// Expects JWT token in query parameter: GET /api/ws?token=<jwt>
func (h *Handler) ServeWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		// Try to get token from Authorization header as fallback
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing or invalid token",
		})
		return
	}

	// Validate JWT
	claims, err := h.jwtManager.ValidateAccessToken(token)
	if err != nil {
		h.log.Warn("Invalid WebSocket token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	// Create connection handler
	client := NewConnection(h.hub, conn, claims.UserID, "", h.log)

	// Register with hub
	h.hub.register <- client

	h.log.Info("WebSocket connection established",
		zap.String("user_id", claims.UserID),
		zap.String("email", claims.Email),
	)

	// Run connection (will block until connection closes)
	client.Run()
}

// ServeRoom handles WebSocket upgrades for room-based connections at /api/ws/room/:room.
// Expects JWT token in query parameter: GET /api/ws/room/chat-room?token=<jwt>
func (h *Handler) ServeRoom(c *gin.Context) {
	roomID := c.Param("room")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "room ID is required",
		})
		return
	}

	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing or invalid token",
		})
		return
	}

	// Validate JWT
	claims, err := h.jwtManager.ValidateAccessToken(token)
	if err != nil {
		h.log.Warn("Invalid WebSocket token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.log.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	// Create connection handler with room ID
	client := NewConnection(h.hub, conn, claims.UserID, roomID, h.log)

	// Register with hub
	h.hub.register <- client

	h.log.Info("WebSocket room connection established",
		zap.String("user_id", claims.UserID),
		zap.String("room_id", roomID),
	)

	// Run connection
	client.Run()
}

// GetStatus returns WebSocket hub status (for health checks).
func (h *Handler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"websocket": gin.H{
			"connected_clients": h.hub.ClientCount(),
			"status":            "running",
		},
	})
}
