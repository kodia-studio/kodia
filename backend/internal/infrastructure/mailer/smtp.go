package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/wneessen/go-mail"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// SMTPMailer implements ports.Mailer using SMTP.
type SMTPMailer struct {
	client    *mail.Client
	config    *config.MailConfig
	log       *zap.Logger
	basePath  string // Path to email templates
}

// NewSMTPMailer creates a new SMTPMailer.
func NewSMTPMailer(cfg *config.Config, log *zap.Logger) (*SMTPMailer, error) {
	c, err := mail.NewClient(cfg.Mail.Host,
		mail.WithPort(cfg.Mail.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Mail.User),
		mail.WithPassword(cfg.Mail.Password),
		mail.WithTLSPolicy(mail.TLSOpportunistic),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail client: %w", err)
	}

	// Default template path
	basePath := "resources/mail"

	return &SMTPMailer{
		client:   c,
		config:   &cfg.Mail,
		log:      log,
		basePath: basePath,
	}, nil
}

func (m *SMTPMailer) Send(ctx context.Context, to []string, subject string, body string) error {
	msg := mail.NewMsg()
	if err := msg.From(fmt.Sprintf("%s <%s>", m.config.FromName, m.config.FromAddr)); err != nil {
		return err
	}
	if err := msg.To(to...); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	return m.client.DialAndSendWithContext(ctx, msg)
}

func (m *SMTPMailer) SendHTML(ctx context.Context, to []string, subject string, htmlBody string) error {
	msg := mail.NewMsg()
	if err := msg.From(fmt.Sprintf("%s <%s>", m.config.FromName, m.config.FromAddr)); err != nil {
		return err
	}
	if err := msg.To(to...); err != nil {
		return err
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextHTML, htmlBody)

	return m.client.DialAndSendWithContext(ctx, msg)
}

func (m *SMTPMailer) SendWithTemplate(ctx context.Context, to []string, subject string, templatePath string, data interface{}) error {
	fullPath := filepath.Join(m.basePath, templatePath)
	
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		m.log.Error("Failed to parse email template", zap.String("path", fullPath), zap.Error(err))
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		m.log.Error("Failed to execute email template", zap.String("path", fullPath), zap.Error(err))
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return m.SendHTML(ctx, to, subject, buf.String())
}
