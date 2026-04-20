package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

// Hub maintains the set of active WebSocket connections and broadcasts messages to them.
type Hub struct {
	// Registered clients
	clients map[*Connection]bool

	// Room-based connections: room name -> set of connections
	rooms map[string]map[*Connection]bool

	// User-based connections: user ID -> slice of connections
	userConns map[string][]*Connection

	// Inbound messages from clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Connection

	// Unregister requests from clients
	unregister chan *Connection

	mu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Connection),
		unregister: make(chan *Connection),
		clients:    make(map[*Connection]bool),
		rooms:      make(map[string]map[*Connection]bool),
		userConns:  make(map[string][]*Connection),
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

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

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
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, skip message
				}
			}
			h.mu.RUnlock()

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
