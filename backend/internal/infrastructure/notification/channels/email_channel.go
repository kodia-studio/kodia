package channels

import (
	"context"

	"github.com/kodia-studio/kodia/internal/core/ports"
)

// EmailChannel delivers notifications via the existing Mailer port.
type EmailChannel struct {
	mailer ports.Mailer
}

// NewEmailChannel creates a new EmailChannel.
func NewEmailChannel(mailer ports.Mailer) *EmailChannel {
	return &EmailChannel{mailer: mailer}
}

func (c *EmailChannel) Name() string { return "email" }

func (c *EmailChannel) Send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	msg := notification.ToNotification("email", notifiable)
	if msg == nil {
		return nil
	}

	to := []string{notifiable.GetEmail()}

	if msg.HtmlBody != "" {
		return c.mailer.SendHTML(ctx, to, msg.Subject, msg.HtmlBody)
	}
	return c.mailer.Send(ctx, to, msg.Subject, msg.TextBody)
}
