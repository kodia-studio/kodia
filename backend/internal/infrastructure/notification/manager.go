// Package notification provides the unified NotificationManager for Kodia Framework.
package notification

import (
	"context"
	"fmt"
	"sync"

	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// Manager routes notifications to the appropriate registered channels.
type Manager struct {
	channels map[string]ports.NotificationChannel
	mu       sync.RWMutex
	log      *zap.Logger
}

// NewManager creates a new NotificationManager.
func NewManager(log *zap.Logger) *Manager {
	return &Manager{
		channels: make(map[string]ports.NotificationChannel),
		log:      log,
	}
}

// Register adds a notification channel driver to the manager.
func (m *Manager) Register(channel ports.NotificationChannel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.channels[channel.Name()] = channel
}

// Send dispatches the notification to all channels returned by Via().
func (m *Manager) Send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	channels := notification.Via(notifiable)
	return m.send(ctx, notifiable, notification, channels)
}

// SendTo dispatches the notification to the specified subset of channels.
func (m *Manager) SendTo(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification, channelNames ...string) error {
	return m.send(ctx, notifiable, notification, channelNames)
}

func (m *Manager) send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification, channelNames []string) error {
	m.mu.RLock()
	channels := make([]ports.NotificationChannel, 0, len(channelNames))
	for _, name := range channelNames {
		ch, ok := m.channels[name]
		if !ok {
			m.log.Warn("notification channel not registered", zap.String("channel", name))
			continue
		}
		channels = append(channels, ch)
	}
	m.mu.RUnlock()

	var firstErr error
	for _, ch := range channels {
		if err := ch.Send(ctx, notifiable, notification); err != nil {
			m.log.Error("failed to send notification",
				zap.String("channel", ch.Name()),
				zap.String("notifiable_id", notifiable.GetID()),
				zap.Error(err),
			)
			if firstErr == nil {
				firstErr = fmt.Errorf("channel %s: %w", ch.Name(), err)
			}
		}
	}
	return firstErr
}
