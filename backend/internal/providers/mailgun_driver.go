package providers

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunProvider implements MailProvider using Mailgun API
type MailgunProvider struct {
	config *MailConfig
	mg     *mailgun.MailgunImpl
}

// NewMailgunProvider creates a new Mailgun mail provider
func NewMailgunProvider(config *MailConfig) (*MailgunProvider, error) {
	if config.Mailgun == nil {
		return nil, fmt.Errorf("Mailgun config is required")
	}

	cfg := config.Mailgun

	// Create Mailgun client
	var mg *mailgun.MailgunImpl
	if cfg.Region == "eu" {
		mg = mailgun.NewMailgun(cfg.Domain, cfg.APIKey)
		mg.SetAPIBase(mailgun.APIBaseEU)
	} else {
		mg = mailgun.NewMailgun(cfg.Domain, cfg.APIKey)
	}

	return &MailgunProvider{
		config: config,
		mg:     mg,
	}, nil
}

// Send sends a single email via Mailgun
func (p *MailgunProvider) Send(ctx context.Context, mail *Mail) error {
	if err := mail.Validate(); err != nil {
		return err
	}

	// Build Mailgun message
	from := mail.BuildFrom(p.config.From, p.config.FromName)
	message := p.mg.NewMessage(from, mail.Subject, mail.Body, mail.To...)

	// Add CC recipients
	for _, cc := range mail.Cc {
		message.AddCC(cc)
	}

	// Add BCC recipients
	for _, bcc := range mail.Bcc {
		message.AddBCC(bcc)
	}

	// Set HTML body if provided
	if mail.HTMLBody != "" {
		message.SetHtml(mail.HTMLBody)
	}

	// Set reply-to if provided
	if mail.ReplyTo != "" {
		message.SetReplyTo(mail.ReplyTo)
	}

	// Add custom headers
	for key, value := range mail.Headers {
		message.AddHeader(key, value)
	}

	// Add metadata as variables
	if mail.Metadata != nil {
		for key, value := range mail.Metadata {
			message.AddVariable(key, fmt.Sprintf("%v", value))
		}
	}

	// Send message
	_, _, err := p.mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via Mailgun: %w", err)
	}

	return nil
}

// SendBatch sends multiple emails via Mailgun
func (p *MailgunProvider) SendBatch(ctx context.Context, mails []*Mail) error {
	for _, mail := range mails {
		if err := p.Send(ctx, mail); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the Mailgun provider
func (p *MailgunProvider) Close() error {
	return nil
}
