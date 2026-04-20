package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 // 512KB
)

// Connection represents a single WebSocket connection.
type Connection struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Message
	UserID string
	RoomID string
	log    *zap.Logger
}

// NewConnection creates a new WebSocket connection handler.
func NewConnection(hub *Hub, conn *websocket.Conn, userID string, roomID string, log *zap.Logger) *Connection {
	return &Connection{
		hub:    hub,
		conn:   conn,
		send:   make(chan *Message, 256),
		UserID: userID,
		RoomID: roomID,
		log:    log,
	}
}

// Run starts the connection's read and write pumps (must be called in a goroutine).
func (c *Connection) Run() {
	go c.writePump()
	c.readPump()
}

// readPump reads messages from the WebSocket connection.
func (c *Connection) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error("WebSocket error", zap.Error(err), zap.String("user_id", c.UserID))
			}
			break
		}

		// Set sender information
		msg.UserID = c.UserID
		msg.RoomID = c.RoomID
		msg.Timestamp = time.Now().Unix()

		// Route message based on type
		switch msg.Type {
		case MessageTypeChat:
			// Send to room if applicable
			if c.RoomID != "" {
				c.hub.SendToRoom(c.RoomID, &msg)
			} else {
				c.hub.Broadcast(&msg)
			}

		case MessageTypePing:
			// Respond with pong
			pongMsg := &Message{
				Type:      MessageTypePong,
				UserID:    c.UserID,
				Timestamp: time.Now().Unix(),
			}
			select {
			case c.send <- pongMsg:
			default:
				// Send channel full
			}

		default:
			// Other messages can be broadcast
			c.hub.Broadcast(&msg)
		}
	}
}

// writePump writes messages to the WebSocket connection.
func (c *Connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the connection
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Marshal and send message
			data, err := json.Marshal(msg)
			if err != nil {
				c.log.Error("Failed to marshal message", zap.Error(err), zap.String("user_id", c.UserID))
				continue
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(data)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				nextMsg := <-c.send
				nextData, _ := json.Marshal(nextMsg)
				w.Write(nextData)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to this connection.
func (c *Connection) SendMessage(msg *Message) {
	select {
	case c.send <- msg:
	default:
		// Send channel is full, skip
	}
}

// Close closes the connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}
