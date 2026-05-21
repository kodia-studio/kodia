package listeners

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/events"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// SendWelcomeNotification sends a welcome notification when a user registers.
type SendWelcomeNotification struct {
	NotificationService ports.NotificationService
	Log                 *zap.Logger
}

func (l *SendWelcomeNotification) Handle(ctx context.Context, event ports.Event) error {
	e, ok := event.Payload().(events.UserRegisteredEvent)
	if !ok {
		return nil
	}

	_, err := l.NotificationService.Send(ctx, ports.SendNotificationInput{
		UserID:    e.UserID,
		Type:      domain.NotificationTypeSuccess,
		Title:     "Welcome to Kodia!",
		Message:   "Your account has been successfully created. Welcome to the Kodia Framework community!",
		SendEmail: false,
	})
	if err != nil {
		l.Log.Warn("Failed to send welcome notification", zap.String("user_id", e.UserID), zap.Error(err))
	}
	return nil
}
