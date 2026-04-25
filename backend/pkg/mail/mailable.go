package mail

import (
	"context"
	"fmt"

	"github.com/kodia-studio/kodia/internal/providers"
)

// Mailable defines the interface for mail classes
type Mailable interface {
	// GetMail returns the Mail message
	GetMail() *providers.Mail

	// GetSubject returns the email subject
	GetSubject() string

	// GetTemplate returns the template path
	GetTemplate() string

	// Build constructs the email content
	Build(ctx context.Context, data map[string]any) error
}

// BaseMail provides a foundation for mail classes
type BaseMail struct {
	Mail *providers.Mail
}

// NewBaseMail creates a new BaseMail instance
func NewBaseMail() *BaseMail {
	return &BaseMail{
		Mail: &providers.Mail{
			Headers:  make(map[string]string),
			Metadata: make(map[string]any),
		},
	}
}

// To sets the recipient(s)
func (m *BaseMail) To(addresses ...string) *BaseMail {
	m.Mail.To = addresses
	return m
}

// AddTo adds additional recipient(s)
func (m *BaseMail) AddTo(addresses ...string) *BaseMail {
	m.Mail.To = append(m.Mail.To, addresses...)
	return m
}

// Cc sets the CC recipient(s)
func (m *BaseMail) Cc(addresses ...string) *BaseMail {
	m.Mail.Cc = addresses
	return m
}

// Bcc sets the BCC recipient(s)
func (m *BaseMail) Bcc(addresses ...string) *BaseMail {
	m.Mail.Bcc = addresses
	return m
}

// From sets the sender address
func (m *BaseMail) From(address string) *BaseMail {
	m.Mail.From = address
	return m
}

// Subject sets the email subject
func (m *BaseMail) Subject(subject string) *BaseMail {
	m.Mail.Subject = subject
	return m
}

// Body sets the plain text body
func (m *BaseMail) Body(content string) *BaseMail {
	m.Mail.Body = content
	return m
}

// HTMLBody sets the HTML body
func (m *BaseMail) HTMLBody(content string) *BaseMail {
	m.Mail.HTMLBody = content
	return m
}

// ReplyTo sets the reply-to address
func (m *BaseMail) ReplyTo(address string) *BaseMail {
	m.Mail.ReplyTo = address
	return m
}

// AddHeader adds a custom header
func (m *BaseMail) AddHeader(key, value string) *BaseMail {
	if m.Mail.Headers == nil {
		m.Mail.Headers = make(map[string]string)
	}
	m.Mail.Headers[key] = value
	return m
}

// AddMetadata adds metadata key-value pair
func (m *BaseMail) AddMetadata(key string, value any) *BaseMail {
	if m.Mail.Metadata == nil {
		m.Mail.Metadata = make(map[string]any)
	}
	m.Mail.Metadata[key] = value
	return m
}

// Attach adds an attachment
func (m *BaseMail) Attach(filename, contentType string, content []byte) *BaseMail {
	attachment := &providers.Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
	}
	m.Mail.Attachments = append(m.Mail.Attachments, attachment)
	return m
}

// GetMail returns the Mail message
func (m *BaseMail) GetMail() *providers.Mail {
	return m.Mail
}

// GetSubject returns the email subject
func (m *BaseMail) GetSubject() string {
	return m.Mail.Subject
}

// GetTemplate returns the template path
func (m *BaseMail) GetTemplate() string {
	return ""
}

// Build constructs the email content (empty implementation)
func (m *BaseMail) Build(ctx context.Context, data map[string]any) error {
	return nil
}

// Sender handles sending mails
type Sender struct {
	provider providers.MailProvider
	queue    *MailQueue
	config   *SenderConfig
}

// SenderConfig contains sender configuration
type SenderConfig struct {
	Provider providers.MailProvider
	Queue    *MailQueue
	UseQueue bool
}

// NewSender creates a new mail sender
func NewSender(provider providers.MailProvider, useQueue bool, queue *MailQueue) *Sender {
	return &Sender{
		provider: provider,
		queue:    queue,
		config: &SenderConfig{
			Provider: provider,
			Queue:    queue,
			UseQueue: useQueue,
		},
	}
}

// Send sends a mail message
func (s *Sender) Send(ctx context.Context, mailable Mailable) error {
	if err := mailable.Build(ctx, nil); err != nil {
		return fmt.Errorf("failed to build mail: %w", err)
	}

	mail := mailable.GetMail()
	if err := mail.Validate(); err != nil {
		return fmt.Errorf("failed to validate mail: %w", err)
	}

	if s.config.UseQueue && s.queue != nil {
		return s.queue.Enqueue(ctx, mail)
	}

	return s.provider.Send(ctx, mail)
}

// SendBatch sends multiple mail messages
func (s *Sender) SendBatch(ctx context.Context, mailables []Mailable) error {
	var mails []*providers.Mail

	for _, mailable := range mailables {
		if err := mailable.Build(ctx, nil); err != nil {
			return fmt.Errorf("failed to build mail: %w", err)
		}

		mail := mailable.GetMail()
		if err := mail.Validate(); err != nil {
			return fmt.Errorf("failed to validate mail: %w", err)
		}

		mails = append(mails, mail)
	}

	if s.config.UseQueue && s.queue != nil {
		for _, mail := range mails {
			if err := s.queue.Enqueue(ctx, mail); err != nil {
				return err
			}
		}
		return nil
	}

	return s.provider.SendBatch(ctx, mails)
}
