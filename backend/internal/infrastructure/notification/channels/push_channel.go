package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kodia-studio/kodia/internal/core/ports"
)

const fcmEndpoint = "https://fcm.googleapis.com/fcm/send"

// PushChannel delivers push notifications via Firebase Cloud Messaging (FCM).
type PushChannel struct {
	serverKey string
}

// NewPushChannel creates a new PushChannel using a FCM Legacy Server Key.
func NewPushChannel(serverKey string) *PushChannel {
	return &PushChannel{serverKey: serverKey}
}

func (c *PushChannel) Name() string { return "push" }

type fcmRequest struct {
	To           string            `json:"to"`
	Notification fcmNotification   `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
}

type fcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c *PushChannel) Send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	msg := notification.ToNotification("push", notifiable)
	if msg == nil || msg.PushTitle == "" {
		return nil
	}

	token := notifiable.GetPushToken()
	if token == "" {
		return fmt.Errorf("push channel: notifiable has no push token")
	}

	payload := fcmRequest{
		To: token,
		Notification: fcmNotification{
			Title: msg.PushTitle,
			Body:  msg.PushBody,
		},
		Data: msg.PushData,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("push channel: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fcmEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("push channel: create request: %w", err)
	}
	req.Header.Set("Authorization", "key="+c.serverKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("push channel: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("push channel: FCM returned status %d", resp.StatusCode)
	}
	return nil
}
