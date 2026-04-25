package providers

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// SMTPProvider implements MailProvider using SMTP
type SMTPProvider struct {
	config *MailConfig
	host   string
	port   int
	auth   smtp.Auth
}

// NewSMTPProvider creates a new SMTP mail provider
func NewSMTPProvider(config *MailConfig) (*SMTPProvider, error) {
	if config.SMTP == nil {
		return nil, fmt.Errorf("SMTP config is required")
	}

	cfg := config.SMTP
	var auth smtp.Auth

	// Create authentication if username is provided
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	}

	return &SMTPProvider{
		config: config,
		host:   cfg.Host,
		port:   cfg.Port,
		auth:   auth,
	}, nil
}

// Send sends a single email via SMTP
func (p *SMTPProvider) Send(ctx context.Context, mail *Mail) error {
	if err := mail.Validate(); err != nil {
		return err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.config.SMTP.Timeout)*time.Second)
	defer cancel()

	// Build email content
	body := p.buildEmailBody(mail)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	c, err := p.dialContext(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer c.Close()

	// Authenticate if needed
	if p.auth != nil {
		if err := c.Auth(p.auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	from := mail.BuildFrom(p.config.From, p.config.FromName)
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	recipients := append(mail.To, mail.Cc...)
	recipients = append(recipients, mail.Bcc...)

	for _, recipient := range recipients {
		if err := c.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get SMTP data writer: %w", err)
	}
	defer w.Close()

	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	return nil
}

// SendBatch sends multiple emails
func (p *SMTPProvider) SendBatch(ctx context.Context, mails []*Mail) error {
	for _, mail := range mails {
		if err := p.Send(ctx, mail); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the SMTP provider
func (p *SMTPProvider) Close() error {
	return nil
}

// buildEmailBody builds the email message body with headers
func (p *SMTPProvider) buildEmailBody(mail *Mail) []byte {
	var buf bytes.Buffer

	// Write standard email headers
	from := mail.BuildFrom(p.config.From, p.config.FromName)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))

	buf.WriteString("To: ")
	for i, to := range mail.To {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(to)
	}
	buf.WriteString("\r\n")

	if len(mail.Cc) > 0 {
		buf.WriteString("Cc: ")
		for i, cc := range mail.Cc {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(cc)
		}
		buf.WriteString("\r\n")
	}

	if mail.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", mail.ReplyTo))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Custom headers
	for key, value := range mail.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Determine content type
	if mail.HTMLBody != "" {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(mail.HTMLBody)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(mail.Body)
	}

	buf.WriteString("\r\n")
	return buf.Bytes()
}

// dialContext dials the SMTP server with context
func (p *SMTPProvider) dialContext(ctx context.Context, addr string) (*smtp.Client, error) {
	// Use net.Dialer with context for timeout
	cfg := p.config.SMTP

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if cfg.TLS {
		if err := client.StartTLS(nil); err != nil {
			client.Close()
			return nil, fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	return client, nil
}
