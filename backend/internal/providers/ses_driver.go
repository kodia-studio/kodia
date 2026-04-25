package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESProvider implements MailProvider using AWS SES
type SESProvider struct {
	config *MailConfig
	client *ses.Client
}

// NewSESProvider creates a new AWS SES mail provider
func NewSESProvider(config *MailConfig) (*SESProvider, error) {
	if config.SES == nil {
		return nil, fmt.Errorf("SES config is required")
	}

	cfg := config.SES

	// Create AWS config
	awsCfg, err := createAWSConfigFromSES(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	// Create SES client
	client := ses.NewFromConfig(awsCfg)

	return &SESProvider{
		config: config,
		client: client,
	}, nil
}

// Send sends a single email via AWS SES
func (p *SESProvider) Send(ctx context.Context, mail *Mail) error {
	if err := mail.Validate(); err != nil {
		return err
	}

	// Build SES request
	from := mail.BuildFrom(p.config.From, p.config.FromName)

	// Build message
	var body types.Body

	if mail.HTMLBody != "" {
		body.Html = &types.Content{
			Data:    aws.String(mail.HTMLBody),
			Charset: aws.String("UTF-8"),
		}
	}

	if mail.Body != "" {
		body.Text = &types.Content{
			Data:    aws.String(mail.Body),
			Charset: aws.String("UTF-8"),
		}
	}

	message := &types.Message{
		Subject: &types.Content{
			Data:    aws.String(mail.Subject),
			Charset: aws.String("UTF-8"),
		},
		Body: &body,
	}

	// Prepare destination
	destination := &types.Destination{
		ToAddresses:  mail.To,
		CcAddresses:  mail.Cc,
		BccAddresses: mail.Bcc,
	}

	// Build SES input
	input := &ses.SendEmailInput{
		Source:      aws.String(from),
		Destination: destination,
		Message:     message,
	}

	// Handle reply-to
	if mail.ReplyTo != "" {
		input.ReplyToAddresses = append(input.ReplyToAddresses, mail.ReplyTo)
	}

	// Handle custom headers
	if len(mail.Headers) > 0 {
		// SES doesn't support custom headers in SendEmail API
		// Use SendRawEmail instead (more complex but supports headers)
		// For now, we'll use SendEmail which ignores custom headers
	}

	// Send email
	result, err := p.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	if result.MessageId == nil || *result.MessageId == "" {
		return fmt.Errorf("failed to send email via SES: no message ID returned")
	}

	return nil
}

// SendBatch sends multiple emails via SES
func (p *SESProvider) SendBatch(ctx context.Context, mails []*Mail) error {
	for _, mail := range mails {
		if err := p.Send(ctx, mail); err != nil {
			// Log error but continue sending other emails
			fmt.Printf("failed to send email to %s: %v\n", strings.Join(mail.To, ", "), err)
		}
	}
	return nil
}

// Close closes the SES provider
func (p *SESProvider) Close() error {
	return nil
}

// createAWSConfigFromSES creates AWS SDK configuration from SES config
func createAWSConfigFromSES(cfg *SESConfig) (aws.Config, error) {
	// Create credentials provider
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	)

	// Load default config with credentials
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return awsCfg, nil
}
