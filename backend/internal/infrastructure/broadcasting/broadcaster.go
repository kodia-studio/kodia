// Package broadcasting bridges the domain Event system with real-time delivery
// via WebSocket and SSE — analogous to Laravel's Broadcasting system.
package broadcasting

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kodia-studio/kodia/internal/adapters/sse"
	ws "github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// EventBroadcaster implements ports.Broadcaster.
// It delivers BroadcastEvents to both WebSocket Hub and SSE Manager.
type EventBroadcaster struct {
	hub     *ws.Hub
	sse     *sse.Manager
	log     *zap.Logger
}

// NewEventBroadcaster creates a new EventBroadcaster.
func NewEventBroadcaster(hub *ws.Hub, sseManager *sse.Manager, log *zap.Logger) *EventBroadcaster {
	return &EventBroadcaster{hub: hub, sse: sseManager, log: log}
}

// Broadcast pushes a BroadcastEvent to all channels it declares via BroadcastOn().
func (b *EventBroadcaster) Broadcast(ctx context.Context, event ports.BroadcastEvent) error {
	eventName := event.BroadcastAs()
	if eventName == "" {
		eventName = event.Name()
	}

	data := event.BroadcastWith()
	if data == nil {
		if p, ok := event.Payload().(map[string]interface{}); ok {
			data = p
		} else {
			data = map[string]interface{}{"payload": event.Payload()}
		}
	}

	channels := event.BroadcastOn()
	for _, channel := range channels {
		b.pushToChannel(channel, eventName, data)
	}
	return nil
}

// BroadcastToUser pushes an ad-hoc event to a specific user's private channels.
func (b *EventBroadcaster) BroadcastToUser(_ context.Context, userID, eventName string, data map[string]interface{}) error {
	b.pushToWS("private-"+userID, eventName, data)
	b.sse.PublishToUser(userID, eventName, data)
	return nil
}

// BroadcastToRoom pushes an ad-hoc event to all clients in a room/presence channel.
func (b *EventBroadcaster) BroadcastToRoom(_ context.Context, roomID, eventName string, data map[string]interface{}) error {
	channel := fmt.Sprintf("presence-%s", roomID)
	b.pushToChannel(channel, eventName, data)
	return nil
}

// pushToChannel routes a message to both WebSocket and SSE based on the channel prefix.
func (b *EventBroadcaster) pushToChannel(channel, eventName string, data map[string]interface{}) {
	b.log.Debug("broadcasting event",
		zap.String("channel", channel),
		zap.String("event", eventName),
	)

	// WebSocket delivery
	b.pushToWS(channel, eventName, data)

	// SSE delivery
	if err := b.sse.Publish(channel, eventName, data); err != nil {
		b.log.Error("sse publish failed", zap.String("channel", channel), zap.Error(err))
	}

	// For presence channels, also push to room via Hub
	if strings.HasPrefix(channel, "presence-") {
		roomID := strings.TrimPrefix(channel, "presence-")
		b.hub.SendToRoom(roomID, &ws.Message{
			Type:      ws.MessageTypeBroadcast,
			Event:     eventName,
			Channel:   channel,
			Payload:   ws.BroadcastPayload{Event: eventName, Channel: channel, Data: data},
			Timestamp: time.Now().Unix(),
		})
	}
}

func (b *EventBroadcaster) pushToWS(channel, eventName string, data map[string]interface{}) {
	msg := &ws.Message{
		Type:      ws.MessageTypeBroadcast,
		Event:     eventName,
		Channel:   channel,
		Payload:   ws.BroadcastPayload{Event: eventName, Channel: channel, Data: data},
		Timestamp: time.Now().Unix(),
	}

	ch := b.hub.GetOrCreateChannel(channel)
	ch.Broadcast(msg)
}
