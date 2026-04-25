# Mailer Provider

Kodia Framework provides a unified email system supporting multiple mail drivers (SMTP, Mailgun, Resend, AWS SES) with template support and queue-based async sending.

---

## Configuration

### Environment Variables

```bash
# Mail Provider
MAIL_DRIVER=smtp              # smtp, mailgun, resend, ses
MAIL_FROM=noreply@example.com
MAIL_FROM_NAME=Kodia
MAIL_QUEUE=true               # Enable queue-based sending

# SMTP Configuration
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USERNAME=your-username
MAIL_PASSWORD=your-password
MAIL_SECURE=true
MAIL_TLS=true
MAIL_TIMEOUT=10

# Mailgun Configuration
MAILGUN_DOMAIN=mail.example.com
MAILGUN_API_KEY=your-api-key
MAILGUN_REGION=us             # us or eu

# Resend Configuration
RESEND_API_KEY=your-api-key

# AWS SES Configuration
SES_REGION=us-east-1
SES_ACCESS_KEY=your-access-key
SES_SECRET_KEY=your-secret-key

# Redis (for queue)
REDIS_HOST=localhost
REDIS_PORT=6379
```

---

## Drivers

### SMTP

Best for: Self-hosted email servers, Gmail, Office 365, custom SMTP

```bash
MAIL_DRIVER=smtp
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
```

### Mailgun

Best for: Reliable email delivery, detailed analytics, high volume

```bash
MAIL_DRIVER=mailgun
MAILGUN_DOMAIN=mail.example.com
MAILGUN_API_KEY=key-xxxxx
MAILGUN_REGION=us
```

**Features**:
- Email tracking
- Deliverability analytics
- Spam filtering
- Custom headers support

### Resend

Best for: Modern SaaS, simple setup, developer-friendly

```bash
MAIL_DRIVER=resend
RESEND_API_KEY=re_xxxxx
```

**Features**:
- Built for developers
- Easy integration
- Good deliverability
- Email templates

### AWS SES

Best for: High volume, AWS infrastructure, cost-effective

```bash
MAIL_DRIVER=ses
SES_REGION=us-east-1
SES_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
SES_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

**Features**:
- High sending limits
- Reputation management
- Production access required
- IAM credentials needed

---

## Basic Usage

### Sending a Simple Email

```go
package handlers

import (
    "context"
    "github.com/kodia-studio/kodia/backend/internal/providers"
    "github.com/kodia-studio/kodia/backend/pkg/mail"
)

func SendEmail(ctx context.Context, sender *mail.Sender) error {
    // Create a simple mail
    m := &providers.Mail{
        From:     "noreply@example.com",
        To:       []string{"user@example.com"},
        Subject:  "Hello!",
        Body:     "This is a simple email.",
        HTMLBody: "<h1>Hello!</h1><p>This is a simple email.</p>",
    }

    return sender.Send(ctx, m)
}
```

### Creating Mail Classes

Create a mail class in `pkg/mail/mails/`:

```go
package mails

import (
    "context"
    "github.com/kodia-studio/kodia/backend/pkg/mail"
)

type NotificationMail struct {
    *mail.BaseMail
    Title   string
    Message string
}

func NewNotificationMail() *NotificationMail {
    return &NotificationMail{
        BaseMail: mail.NewBaseMail(),
    }
}

func (nm *NotificationMail) SetTitle(title string) *NotificationMail {
    nm.Title = title
    nm.Subject(title)
    return nm
}

func (nm *NotificationMail) SetMessage(message string) *NotificationMail {
    nm.Message = message
    nm.HTMLBody("<p>" + message + "</p>")
    return nm
}

func (nm *NotificationMail) Build(ctx context.Context, data map[string]any) error {
    // Custom build logic
    return nil
}
```

---

## Template-Based Emails

### Using Templates

Templates are stored in `resources/mail-templates/`:

```
resources/
└── mail-templates/
    ├── welcome.html
    ├── welcome.txt
    ├── verification.html
    └── verification.txt
```

### Creating a Template Mail

```go
package mails

import (
    "context"
    "github.com/kodia-studio/kodia/backend/pkg/mail"
)

type WelcomeMail struct {
    *mail.TemplateMail
    UserName string
    Email    string
}

func NewWelcomeMail(engine *mail.TemplateEngine) *WelcomeMail {
    return &WelcomeMail{
        TemplateMail: mail.NewTemplateMail(engine),
    }
}

func (wm *WelcomeMail) SetEmail(email string) *WelcomeMail {
    wm.Email = email
    wm.To(email)
    return wm
}

func (wm *WelcomeMail) SetUserName(name string) *WelcomeMail {
    wm.UserName = name
    return wm
}

func (wm *WelcomeMail) Build(ctx context.Context, data map[string]any) error {
    wm.Subject("Welcome to Kodia!")
    wm.WithTemplate("welcome")
    wm.SetTemplateVariable("name", wm.UserName)
    wm.SetTemplateVariable("email", wm.Email)
    return wm.TemplateMail.Build(ctx, data)
}
```

### Template File (HTML)

```html
<!DOCTYPE html>
<html>
<head>
    <title>Welcome</title>
</head>
<body>
    <h1>Welcome, {{ .name }}!</h1>
    <p>Thank you for signing up at {{ .app_name }}.</p>
    <p>Your email: {{ .email }}</p>
</body>
</html>
```

---

## Queue-Based Sending

### Configuration

Enable queue-based sending:

```bash
MAIL_QUEUE=true
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Sending to Queue

```go
// Mail is automatically queued based on MAIL_QUEUE config
sender := mail.NewSender(provider, true, mailQueue)

mail := mails.NewWelcomeMail(templateEngine)
mail.SetEmail("user@example.com").SetUserName("John")

// This will queue the email instead of sending immediately
sender.Send(ctx, mail)
```

### Processing Queue

Start the mail queue worker:

```go
package main

import (
    "github.com/kodia-studio/kodia/backend/pkg/mail"
)

func main() {
    // Start mail queue processor
    if err := mail.ProcessMailQueue("localhost:6379", provider, 10); err != nil {
        panic(err)
    }
}
```

### Delayed Sending

```go
// Send email after 1 hour
mailQueue.EnqueueDelayed(ctx, mail, 1*time.Hour)

// Or with custom options
import "github.com/asynq/asynq"

mailQueue.EnqueueWithOptions(ctx, mail,
    asynq.ProcessIn(1*time.Hour),
    asynq.MaxRetry(3),
)
```

---

## Advanced Usage

### Custom Headers

```go
mail := mails.NewWelcomeMail(templateEngine)
mail.To("user@example.com")
mail.AddHeader("X-Mailer", "Kodia/1.7.0")
mail.AddHeader("X-Priority", "1")
```

### Metadata

```go
mail := mails.NewWelcomeMail(templateEngine)
mail.AddMetadata("user_id", "123")
mail.AddMetadata("campaign", "welcome")
```

### Multiple Recipients

```go
mail := mails.NewWelcomeMail(templateEngine)
mail.To("user1@example.com", "user2@example.com")
mail.Cc("manager@example.com")
mail.Bcc("archive@example.com")
```

### Attachments

```go
mail := mails.NewWelcomeMail(templateEngine)
mail.Attach("invoice.pdf", "application/pdf", pdfData)
mail.Attach("logo.png", "image/png", imageData)
```

### Custom Template Data

```go
engine := mail.NewTemplateEngine("resources/mail-templates")

mail := mails.NewWelcomeMail(engine)
mail.To("user@example.com")
mail.SetTemplateVariable("company_name", "ACME Corp")
mail.SetTemplateVariable("support_email", "help@acme.com")
```

---

## Examples

### User Registration

```go
func RegisterUser(ctx context.Context, user *User, sender *mail.Sender, engine *mail.TemplateEngine) error {
    // Create welcome mail
    welcomeMail := mails.NewWelcomeMail(engine)
    welcomeMail.SetEmail(user.Email).SetUserName(user.Name)

    if err := sender.Send(ctx, welcomeMail); err != nil {
        return err
    }

    // Create verification mail
    verificationMail := mails.NewVerificationMail(engine)
    verificationMail.
        SetEmail(user.Email).
        SetUserName(user.Name).
        SetVerifyLink("https://example.com/verify?token=123")

    return sender.Send(ctx, verificationMail)
}
```

### Batch Sending

```go
func SendNewsletterTo(ctx context.Context, users []*User, sender *mail.Sender, engine *mail.TemplateEngine) error {
    var mailables []mail.Mailable

    for _, user := range users {
        m := mails.NewNewsletterMail(engine)
        m.SetEmail(user.Email).SetUserName(user.Name)
        mailables = append(mailables, m)
    }

    return sender.SendBatch(ctx, mailables)
}
```

---

## Testing

### Mock Provider

```go
type MockMailProvider struct {
    SentMails []*providers.Mail
}

func (m *MockMailProvider) Send(ctx context.Context, mail *providers.Mail) error {
    m.SentMails = append(m.SentMails, mail)
    return nil
}
```

### Test Example

```go
func TestWelcomeMail(t *testing.T) {
    mockProvider := &MockMailProvider{}
    sender := mail.NewSender(mockProvider, false, nil)
    engine := mail.NewTemplateEngine("resources/mail-templates")

    m := mails.NewWelcomeMail(engine)
    m.SetEmail("test@example.com").SetUserName("Test User")

    err := sender.Send(context.Background(), m)
    if err != nil {
        t.Fatalf("Send failed: %v", err)
    }

    if len(mockProvider.SentMails) != 1 {
        t.Errorf("Expected 1 mail, got %d", len(mockProvider.SentMails))
    }
}
```

---

## Best Practices

✅ **Do**:
- Use templates for consistent branding
- Enable queue for better UX
- Set appropriate retry policies
- Use metadata for tracking
- Test locally with SMTP
- Monitor delivery rates
- Handle errors gracefully

❌ **Don't**:
- Send passwords in emails
- Use plain text for sensitive data
- Hardcode email addresses
- Ignore bounces/failures
- Send without validation
- Expose API keys in config files
- Send too frequently to same user

---

## Troubleshooting

### Emails Not Sending

1. Check configuration
2. Verify credentials
3. Check Redis for queue issues
4. Review mail logs
5. Test with mock provider

### High Bounce Rate

- Verify email addresses
- Check sender reputation
- Review email content
- Test with different provider

### Slow Delivery

- Enable queue-based sending
- Increase worker concurrency
- Check Redis performance
- Review provider limits

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
