package websocket

// Broadcaster provides a service-layer interface for sending WebSocket messages.
// Inject this into your services to trigger real-time updates.
type Broadcaster struct {
	hub *Hub
}

// NewBroadcaster creates a new Broadcaster instance.
func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{hub: hub}
}

// NotifyUser sends a notification to a specific user.
// Example:
//  broadcaster.NotifyUser(userID, "order-update", OrderNotificationPayload{
//      OrderID: "123",
//      Status: "shipped",
//  })
func (b *Broadcaster) NotifyUser(userID string, eventType string, payload interface{}) error {
	msg := &Message{
		Type:    MessageTypeNotification,
		Payload: payload,
		UserID:  userID,
	}
	b.hub.SendToUser(userID, msg)
	return nil
}

// NotifyMultipleUsers sends a notification to multiple users.
func (b *Broadcaster) NotifyMultipleUsers(userIDs []string, eventType string, payload interface{}) error {
	for _, userID := range userIDs {
		b.NotifyUser(userID, eventType, payload)
	}
	return nil
}

// Broadcast sends a message to all connected clients.
func (b *Broadcaster) Broadcast(eventType string, payload interface{}) error {
	msg := &Message{
		Type:    MessageTypeNotification,
		Payload: payload,
	}
	b.hub.Broadcast(msg)
	return nil
}

// BroadcastToRoom sends a message to all clients in a specific room.
// Example:
//  broadcaster.BroadcastToRoom("chat-room-1", "message", ChatPayload{
//      SenderID: "user-123",
//      Content: "Hello everyone!",
//  })
func (b *Broadcaster) BroadcastToRoom(roomID string, eventType string, payload interface{}) error {
	msg := &Message{
		Type:    MessageTypeChat,
		Payload: payload,
		RoomID:  roomID,
	}
	b.hub.SendToRoom(roomID, msg)
	return nil
}

// SendError sends an error message to a specific user.
func (b *Broadcaster) SendError(userID string, code string, message string) error {
	msg := &Message{
		Type: MessageTypeError,
		Payload: ErrorPayload{
			Code:    code,
			Message: message,
		},
		UserID: userID,
	}
	b.hub.SendToUser(userID, msg)
	return nil
}

// SendStatus sends a status update to all connected clients.
func (b *Broadcaster) SendStatus(message string) error {
	msg := &Message{
		Type: MessageTypeStatus,
		Payload: StatusPayload{
			Message: message,
		},
	}
	b.hub.Broadcast(msg)
	return nil
}

// ClientCount returns the total number of connected clients.
func (b *Broadcaster) ClientCount() int {
	return b.hub.ClientCount()
}

// UserConnCount returns the number of connections for a specific user.
func (b *Broadcaster) UserConnCount(userID string) int {
	return b.hub.UserConnCount(userID)
}

// RoomConnCount returns the number of connections in a specific room.
func (b *Broadcaster) RoomConnCount(roomID string) int {
	return b.hub.RoomConnCount(roomID)
}
