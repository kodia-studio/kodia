package ports

import (
	"context"
)

// Mailer defines the interface for sending emails.
type Mailer interface {
	// Send sends a plain text email.
	Send(ctx context.Context, to []string, subject string, body string) error
	// SendHTML sends an HTML email.
	SendHTML(ctx context.Context, to []string, subject string, htmlBody string) error
	// SendWithTemplate renders a template and sends it.
	// templatePath is relative to the mail resources directory.
	SendWithTemplate(ctx context.Context, to []string, subject string, templatePath string, data interface{}) error
}
