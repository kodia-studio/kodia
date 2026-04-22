// Package websocket provides WebSocket server implementation for Kodia Framework.
package websocket

// MessageType represents the type of WebSocket message.
type MessageType string

const (
	MessageTypeNotification MessageType = "notification"
	MessageTypeChat         MessageType = "chat"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
	MessageTypeError        MessageType = "error"
	MessageTypeStatus       MessageType = "status"
	// MessageTypeBroadcast is sent when a domain BroadcastEvent is pushed to clients.
	MessageTypeBroadcast MessageType = "broadcast"
	// MessageTypePresence is sent when a user joins or leaves a presence channel.
	MessageTypePresence MessageType = "presence"
)

// Message represents a WebSocket message with typed payload.
type Message struct {
	Type      MessageType `json:"type"`
	Event     string      `json:"event,omitempty"` // broadcast event name
	Channel   string      `json:"channel,omitempty"`
	Payload   interface{} `json:"payload"`
	RoomID    string      `json:"room_id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// NotificationPayload represents a user notification.
type NotificationPayload struct {
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ChatPayload represents a chat message.
type ChatPayload struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
	RoomID     string `json:"room_id"`
}

// ErrorPayload represents an error message.
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// StatusPayload represents a status update.
type StatusPayload struct {
	Message string `json:"message"`
}

// BroadcastPayload carries a domain event pushed via the broadcasting system.
type BroadcastPayload struct {
	Event   string                 `json:"event"`
	Channel string                 `json:"channel"`
	Data    map[string]interface{} `json:"data"`
}

// PresencePayload notifies clients about presence changes in a channel.
type PresencePayload struct {
	Channel string   `json:"channel"`
	UserID  string   `json:"user_id"`
	Joined  bool     `json:"joined"` // true = joined, false = left
	Members []string `json:"members"`
}

