package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kodia-studio/kodia/internal/core/ports"
)

// SlackChannel delivers notifications via a Slack Incoming Webhook.
type SlackChannel struct {
	webhookURL string
}

// NewSlackChannel creates a new SlackChannel.
// webhookURL is the Slack Incoming Webhook URL from your Slack App settings.
func NewSlackChannel(webhookURL string) *SlackChannel {
	return &SlackChannel{webhookURL: webhookURL}
}

func (c *SlackChannel) Name() string { return "slack" }

type slackPayload struct {
	Text        string             `json:"text,omitempty"`
	Attachments []slackAttachment  `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Color string `json:"color,omitempty"`
	Text  string `json:"text"`
}

func (c *SlackChannel) Send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	msg := notification.ToNotification("slack", notifiable)
	if msg == nil || msg.SlackText == "" {
		return nil
	}

	payload := slackPayload{}
	if msg.SlackColor != "" {
		payload.Attachments = []slackAttachment{
			{Color: msg.SlackColor, Text: msg.SlackText},
		}
	} else {
		payload.Text = msg.SlackText
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack channel: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack channel: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("slack channel: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack channel: unexpected status %d", resp.StatusCode)
	}
	return nil
}
