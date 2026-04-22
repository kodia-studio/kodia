package listeners

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/domain"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

// SendNotificationEmail sends an email notification to the user.
type SendNotificationEmail struct {
	Mailer   ports.Mailer
	UserRepo ports.UserRepository
}

// Handle implements ports.Listener interface.
func (l *SendNotificationEmail) Handle(ctx context.Context, event ports.Event) error {
	n, ok := event.Payload().(*domain.Notification)
	if !ok {
		return nil
	}

	user, err := l.UserRepo.FindByID(ctx, n.UserID)
	if err != nil {
		return err
	}

	return l.Mailer.SendWithTemplate(ctx, []string{user.Email}, n.Title, "notification", map[string]interface{}{
		"title":   n.Title,
		"message": n.Message,
		"data":    n.Data,
	})
}
