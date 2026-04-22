package channels

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kodia-studio/kodia/internal/core/ports"
)

// SMSChannel delivers notifications via Twilio SMS API.
type SMSChannel struct {
	accountSID string
	authToken  string
	fromNumber string
}

// NewSMSChannel creates a new SMSChannel configured for Twilio.
func NewSMSChannel(accountSID, authToken, fromNumber string) *SMSChannel {
	return &SMSChannel{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

func (c *SMSChannel) Name() string { return "sms" }

func (c *SMSChannel) Send(ctx context.Context, notifiable ports.Notifiable, notification ports.Notification) error {
	msg := notification.ToNotification("sms", notifiable)
	if msg == nil || msg.SMSText == "" {
		return nil
	}

	to := notifiable.GetPhoneNumber()
	if to == "" {
		return fmt.Errorf("sms channel: notifiable has no phone number")
	}

	return c.sendViaTwilio(ctx, to, msg.SMSText)
}

func (c *SMSChannel) sendViaTwilio(ctx context.Context, to, body string) error {
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.accountSID)

	data := strings.NewReader(fmt.Sprintf("To=%s&From=%s&Body=%s", to, c.fromNumber, body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		return fmt.Errorf("sms channel: create request: %w", err)
	}
	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sms channel: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var twilioErr struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}
		json.NewDecoder(resp.Body).Decode(&twilioErr)
		return fmt.Errorf("sms channel: twilio error %d: %s", twilioErr.Code, twilioErr.Message)
	}
	return nil
}
