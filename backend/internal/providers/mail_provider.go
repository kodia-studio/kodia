package providers

import (
	"context"
	"fmt"
)

// MailProvider defines the interface for sending emails
type MailProvider interface {
	Send(ctx context.Context, mail *Mail) error
	SendBatch(ctx context.Context, mails []*Mail) error
	Close() error
}

// Mail represents an email message
type Mail struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Cc          []string          `json:"cc"`
	Bcc         []string          `json:"bcc"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTMLBody    string            `json:"html_body"`
	ReplyTo     string            `json:"reply_to"`
	Headers     map[string]string `json:"headers"`
	Attachments []*Attachment     `json:"attachments"`
	Metadata    map[string]any    `json:"metadata"`
}

// Attachment represents a file attachment
type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
}

// MailConfig contains configuration for the mail provider
type MailConfig struct {
	Driver   string                 `mapstructure:"driver" default:"smtp"`
	From     string                 `mapstructure:"from"`
	FromName string                 `mapstructure:"from_name" default:"Kodia"`
	SMTP     *SMTPConfig            `mapstructure:"smtp"`
	Mailgun  *MailgunConfig         `mapstructure:"mailgun"`
	Resend   *ResendConfig          `mapstructure:"resend"`
	SES      *SESConfig             `mapstructure:"ses"`
	Queue    bool                   `mapstructure:"queue" default:"true"`
	MaxRetry int                    `mapstructure:"max_retry" default:"3"`
}

// SMTPConfig contains SMTP configuration
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port" default:"587"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Secure   bool   `mapstructure:"secure" default:"true"`
	TLS      bool   `mapstructure:"tls" default:"true"`
	Timeout  int    `mapstructure:"timeout" default:"10"` // seconds
}

// MailgunConfig contains Mailgun API configuration
type MailgunConfig struct {
	Domain string `mapstructure:"domain"`
	APIKey string `mapstructure:"api_key"`
	Region string `mapstructure:"region" default:"us"` // us or eu
}

// ResendConfig contains Resend API configuration
type ResendConfig struct {
	APIKey string `mapstructure:"api_key"`
}

// SESConfig contains AWS SES configuration
type SESConfig struct {
	Region    string `mapstructure:"region" default:"us-east-1"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// NewMailProvider creates a new mail provider based on config
func NewMailProvider(config *MailConfig) (MailProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("mail config is required")
	}

	switch config.Driver {
	case "smtp":
		if config.SMTP == nil {
			return nil, fmt.Errorf("SMTP config is required for smtp driver")
		}
		return NewSMTPProvider(config)

	case "mailgun":
		if config.Mailgun == nil {
			return nil, fmt.Errorf("Mailgun config is required for mailgun driver")
		}
		return NewMailgunProvider(config)

	case "resend":
		if config.Resend == nil {
			return nil, fmt.Errorf("Resend config is required for resend driver")
		}
		return NewResendProvider(config)

	case "ses":
		if config.SES == nil {
			return nil, fmt.Errorf("SES config is required for ses driver")
		}
		return NewSESProvider(config)

	default:
		return nil, fmt.Errorf("unsupported mail driver: %s", config.Driver)
	}
}

// BuildFrom builds the full from address with name
func (m *Mail) BuildFrom(defaultFrom, defaultFromName string) string {
	from := m.From
	if from == "" {
		from = defaultFrom
	}

	if m.Metadata != nil {
		if name, ok := m.Metadata["from_name"]; ok {
			if nameStr, ok := name.(string); ok {
				return fmt.Sprintf("%s <%s>", nameStr, from)
			}
		}
	}

	if defaultFromName != "" {
		return fmt.Sprintf("%s <%s>", defaultFromName, from)
	}

	return from
}

// Validate validates the mail message
func (m *Mail) Validate() error {
	if len(m.To) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}

	if m.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	if m.Body == "" && m.HTMLBody == "" {
		return fmt.Errorf("either body or html_body is required")
	}

	return nil
}
