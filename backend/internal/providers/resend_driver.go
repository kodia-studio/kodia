package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/resend/resend-go/v2"
)

// ResendProvider implements MailProvider using Resend API
type ResendProvider struct {
	config *MailConfig
	client *resend.Client
}

// NewResendProvider creates a new Resend mail provider
func NewResendProvider(config *MailConfig) (*ResendProvider, error) {
	if config.Resend == nil {
		return nil, fmt.Errorf("Resend config is required")
	}

	if config.Resend.APIKey == "" {
		return nil, fmt.Errorf("Resend API key is required")
	}

	return &ResendProvider{
		config: config,
		client: resend.NewClient(config.Resend.APIKey),
	}, nil
}

// Send sends a single email via Resend
func (p *ResendProvider) Send(ctx context.Context, mail *Mail) error {
	if err := mail.Validate(); err != nil {
		return err
	}

	// Build Resend request
	from := mail.BuildFrom(p.config.From, p.config.FromName)

	// Build request parameters
	params := resend.SendEmailRequest{
		From:    from,
		To:      mail.To,
		Subject: mail.Subject,
		Text:    mail.Body,
	}

	// Set HTML body if provided
	if mail.HTMLBody != "" {
		params.Html = mail.HTMLBody
	}

	// Set reply-to if provided
	if mail.ReplyTo != "" {
		params.ReplyTo = mail.ReplyTo
	}

	// Handle CC and BCC
	if len(mail.Cc) > 0 {
		params.Cc = mail.Cc
	}

	if len(mail.Bcc) > 0 {
		params.Bcc = mail.Bcc
	}

	// Handle custom headers (Resend supports limited custom headers)
	if mail.Headers != nil {
		// Initialize headers map if needed
		if params.Headers == nil {
			params.Headers = make(map[string]string)
		}
		// Set headers as a map
		for key, value := range mail.Headers {
			params.Headers[key] = value
		}
	}

	// Send email
	sent, err := p.client.Emails.Send(&params)
	if err != nil {
		return fmt.Errorf("failed to send email via Resend: %w", err)
	}

	if sent.Id == "" {
		return fmt.Errorf("failed to send email via Resend: no message ID returned")
	}

	return nil
}

// SendBatch sends multiple emails via Resend
func (p *ResendProvider) SendBatch(ctx context.Context, mails []*Mail) error {
	for _, mail := range mails {
		if err := p.Send(ctx, mail); err != nil {
			// Log error but continue sending other emails
			fmt.Printf("failed to send email to %s: %v\n", strings.Join(mail.To, ", "), err)
		}
	}
	return nil
}

// Close closes the Resend provider
func (p *ResendProvider) Close() error {
	return nil
}
