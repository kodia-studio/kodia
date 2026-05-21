package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Hub maintains the set of active WebSocket connections and broadcasts messages to them.
type Hub struct {
	// Registered clients
	clients map[*Connection]bool

	// Room-based connections: room name -> set of connections
	rooms map[string]map[*Connection]bool

	// User-based connections: user ID -> slice of connections
	userConns map[string][]*Connection

	// Named channel registry (for broadcasting system)
	channels map[string]*Channel

	// Inbound messages from clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Connection

	// Unregister requests from clients
	unregister chan *Connection

	// Metrics
	totalMessages      atomic.Int64
	totalConnections   atomic.Int64 // Total connections ever made
	totalDisconnections atomic.Int64 // Total disconnections ever made
	peakConnections    atomic.Int64 // Peak concurrent connections
	failedBroadcasts   atomic.Int64 // Messages dropped due to full channels
	startTime          time.Time

	log *zap.Logger

	mu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub(log *zap.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		clients:    make(map[*Connection]bool),
		rooms:      make(map[string]map[*Connection]bool),
		userConns:  make(map[string][]*Connection),
		channels:   make(map[string]*Channel),
		log:        log,
		startTime:  time.Now(),
	}
}

// Run starts the Hub's event loop (must be called in a goroutine).
func (h *Hub) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true

			// Track metrics
			h.totalConnections.Add(1)
			currentConnections := int64(len(h.clients))
			for {
				peak := h.peakConnections.Load()
				if currentConnections <= peak {
					break
				}
				if h.peakConnections.CompareAndSwap(peak, currentConnections) {
					break
				}
			}

			// Register user connection
			if client.UserID != "" {
				h.userConns[client.UserID] = append(h.userConns[client.UserID], client)
			}

			// Register room connection
			if client.RoomID != "" {
				if h.rooms[client.RoomID] == nil {
					h.rooms[client.RoomID] = make(map[*Connection]bool)
				}
				h.rooms[client.RoomID][client] = true
			}

			h.mu.Unlock()
			h.log.Debug("Client connected", zap.String("user_id", client.UserID), zap.Int64("total_connections", currentConnections))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Track metrics
				h.totalDisconnections.Add(1)
				currentConnections := int64(len(h.clients))

				// Unregister user connection
				if client.UserID != "" {
					conns := h.userConns[client.UserID]
					for i, c := range conns {
						if c == client {
							h.userConns[client.UserID] = append(conns[:i], conns[i+1:]...)
							if len(h.userConns[client.UserID]) == 0 {
								delete(h.userConns, client.UserID)
							}
							break
						}
					}
				}

				// Unregister room connection
				if client.RoomID != "" && h.rooms[client.RoomID] != nil {
					delete(h.rooms[client.RoomID], client)
					if len(h.rooms[client.RoomID]) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}

				h.log.Debug("Client disconnected", zap.String("user_id", client.UserID), zap.Int64("total_connections", currentConnections))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.totalMessages.Add(1)
			h.mu.RLock()
			failedCount := 0
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, skip message
					failedCount++
				}
			}
			h.mu.RUnlock()
			if failedCount > 0 {
				h.failedBroadcasts.Add(int64(failedCount))
				h.log.Debug("Broadcast messages dropped", zap.Int("count", failedCount))
			}

		case <-ticker.C:
			// Send periodic ping to all clients
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- &Message{
					Type:      MessageTypePing,
					Timestamp: time.Now().Unix(),
				}:
				default:
					// Skip if client's send channel is full
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients.
func (h *Hub) Broadcast(msg *Message) {
	msg.Timestamp = time.Now().Unix()
	h.broadcast <- msg
}

// SendToUser sends a message to all connections of a specific user.
func (h *Hub) SendToUser(userID string, msg *Message) {
	if userID == "" {
		return
	}

	msg.UserID = userID
	msg.Timestamp = time.Now().Unix()

	h.mu.RLock()
	conns, exists := h.userConns[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	for _, client := range conns {
		select {
		case client.send <- msg:
		default:
			// Client's send channel is full, skip
		}
	}
}

// SendToRoom sends a message to all clients in a specific room.
func (h *Hub) SendToRoom(roomID string, msg *Message) {
	if roomID == "" {
		return
	}

	msg.RoomID = roomID
	msg.Timestamp = time.Now().Unix()

	h.mu.RLock()
	clients, exists := h.rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	for client := range clients {
		select {
		case client.send <- msg:
		default:
			// Client's send channel is full, skip
		}
	}
}

// ClientCount returns the total number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// UserConnCount returns the number of connections for a specific user.
func (h *Hub) UserConnCount(userID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userConns[userID])
}

// RoomConnCount returns the number of connections in a specific room.
func (h *Hub) RoomConnCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms[roomID])
}

// SendJSON is a convenience method to send a JSON-marshaled message.
func (h *Hub) SendJSON(data interface{}) (*Message, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Type:      MessageTypeStatus,
		Payload:   string(jsonBytes),
		Timestamp: time.Now().Unix(),
	}

	h.broadcast <- msg
	return msg, nil
}

// GetOnlineUsers returns the list of user IDs that have at least one active connection.
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users := make([]string, 0, len(h.userConns))
	for uid := range h.userConns {
		users = append(users, uid)
	}
	return users
}

// GetRoomPresence returns the user IDs of users connected to a room.
func (h *Hub) GetRoomPresence(roomID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients, ok := h.rooms[roomID]
	if !ok {
		return nil
	}
	seen := make(map[string]bool)
	members := make([]string, 0)
	for c := range clients {
		if c.UserID != "" && !seen[c.UserID] {
			seen[c.UserID] = true
			members = append(members, c.UserID)
		}
	}
	return members
}

// GetOrCreateChannel returns an existing channel or creates a new one.
func (h *Hub) GetOrCreateChannel(name string) *Channel {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.channels[name]; ok {
		return ch
	}
	ch := NewChannel(name, ParseChannelType(name))
	h.channels[name] = ch
	return ch
}

// BroadcastToChannel delivers a message to a named channel.
func (h *Hub) BroadcastToChannel(channelName string, msg *Message) {
	h.mu.RLock()
	ch, ok := h.channels[channelName]
	h.mu.RUnlock()
	if !ok {
		return
	}
	msg.Channel = channelName
	msg.Timestamp = time.Now().Unix()
	ch.Broadcast(msg)
}

// ActiveRooms returns a list of room names that currently have connections.
func (h *Hub) ActiveRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	rooms := make([]string, 0, len(h.rooms))
	for r := range h.rooms {
		rooms = append(rooms, r)
	}
	return rooms
}

// TotalMessages returns the total number of messages broadcast since startup.
func (h *Hub) TotalMessages() int64 {
	return h.totalMessages.Load()
}

// Metrics returns a snapshot of Hub metrics.
func (h *Hub) Metrics() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	uptime := time.Since(h.startTime)
	totalConns := h.totalConnections.Load()
	totalDisconns := h.totalDisconnections.Load()
	currentConns := int64(len(h.clients))

	return map[string]interface{}{
		// Current state
		"active_connections": currentConns,
		"active_users":       len(h.userConns),
		"active_rooms":       len(h.rooms),
		"active_channels":    len(h.channels),

		// Cumulative metrics
		"total_connections":     totalConns,
		"total_disconnections":  totalDisconns,
		"peak_connections":      h.peakConnections.Load(),
		"total_messages":        h.totalMessages.Load(),
		"failed_broadcasts":     h.failedBroadcasts.Load(),

		// Timing
		"uptime_seconds": uptime.Seconds(),
		"start_time":     h.startTime.Unix(),

		// Calculated
		"reconnections_estimated": totalConns - totalDisconns - currentConns,
	}
}

// --- Broadcaster Interface Implementation ---

// BroadcastToUser implements ports.Broadcaster - sends message to specific user's connections.
func (h *Hub) BroadcastToUser(ctx context.Context, userID string, eventName string, data map[string]interface{}) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	msg := &Message{
		Type:      MessageTypeNotification,
		Event:     eventName,
		Payload:   data,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
	}

	h.SendToUser(userID, msg)
	return nil
}

// BroadcastToRoom implements ports.Broadcaster - sends message to all users in a room.
func (h *Hub) BroadcastToRoom(ctx context.Context, roomID string, eventName string, data map[string]interface{}) error {
	if roomID == "" {
		return fmt.Errorf("roomID cannot be empty")
	}

	msg := &Message{
		Type:      MessageTypeNotification,
		Event:     eventName,
		Payload:   data,
		RoomID:    roomID,
		Timestamp: time.Now().Unix(),
	}

	h.SendToRoom(roomID, msg)
	return nil
}
