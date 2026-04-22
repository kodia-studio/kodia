package ports

import "context"

// BroadcastEvent is implemented by domain events that should be pushed to
// real-time clients via WebSocket or SSE — analogous to Laravel Broadcasting.
type BroadcastEvent interface {
	Event

	// BroadcastOn returns the list of channel names to broadcast on.
	// Examples: "public.orders", "private.user.123", "presence.chat.room1"
	BroadcastOn() []string

	// BroadcastAs overrides the event name sent to the client.
	// Defaults to Event.Name() if empty string is returned.
	BroadcastAs() string

	// BroadcastWith returns the data payload sent to the client.
	// Defaults to Event.Payload() if nil is returned.
	BroadcastWith() map[string]interface{}
}

// Broadcaster delivers BroadcastEvents to real-time clients.
type Broadcaster interface {
	// Broadcast pushes the event to all configured channels (WS + SSE).
	Broadcast(ctx context.Context, event BroadcastEvent) error

	// BroadcastToUser pushes an ad-hoc event to a specific user's private channels.
	BroadcastToUser(ctx context.Context, userID string, eventName string, data map[string]interface{}) error

	// BroadcastToRoom pushes an ad-hoc event to all clients in a room/presence channel.
	BroadcastToRoom(ctx context.Context, roomID string, eventName string, data map[string]interface{}) error
}

// SSEPublisher publishes events to Server-Sent Events clients.
type SSEPublisher interface {
	// Publish sends an event to all SSE clients subscribed to the given channel.
	Publish(channel string, eventName string, data interface{}) error

	// PublishToUser sends an event to a specific user's SSE connection.
	PublishToUser(userID string, eventName string, data interface{}) error
}
