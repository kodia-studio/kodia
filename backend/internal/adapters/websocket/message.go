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
)

// Message represents a WebSocket message with typed payload.
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
	RoomID  string      `json:"room_id,omitempty"`
	UserID  string      `json:"user_id,omitempty"`
	Timestamp int64      `json:"timestamp"`
}

// NotificationPayload represents a user notification.
type NotificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ChatPayload represents a chat message.
type ChatPayload struct {
	SenderID  string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content   string `json:"content"`
	RoomID    string `json:"room_id"`
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
