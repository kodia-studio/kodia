package channels

import (
	"context"
	"time"

	ws "github.com/kodia-studio/kodia/internal/adapters/websocket"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// WebSocketChannel delivers real-time in-app notifications via the WebSocket Hub.
type WebSocketChannel struct {
	hub *ws.Hub
}

// NewWebSocketChannel creates a new WebSocketChannel.
func NewWebSocketChannel(hub *ws.Hub) *WebSocketChannel {
	return &WebSocketChannel{hub: hub}
}

func (c *WebSocketChannel) Name() string { return "websocket" }

func (c *WebSocketChannel) Send(_ context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	msg := notification.ToNotification("websocket", notifiable)
	if msg == nil {
		return nil
	}

	wsMsg := &ws.Message{
		Type:  ws.MessageTypeNotification,
		Event: msg.WSEvent,
		Payload: ws.NotificationPayload{
			Title:   msg.PushTitle,
			Message: msg.TextBody,
		},
		UserID:    notifiable.GetID(),
		Timestamp: time.Now().Unix(),
	}

	if msg.WSPayload != nil {
		wsMsg.Payload = msg.WSPayload
	}

	c.hub.SendToUser(notifiable.GetID(), wsMsg)
	return nil
}
