package ports

import "context"

// Notifiable represents any entity that can receive notifications (e.g., User).
type Notifiable interface {
	// GetID returns the unique identifier of the notifiable entity.
	GetID() string
	// GetEmail returns the email address for email notifications.
	GetEmail() string
	// GetPhoneNumber returns the phone number for SMS notifications.
	GetPhoneNumber() string
	// GetPushToken returns the FCM/APNS push token for push notifications.
	GetPushToken() string
}

// NotificationMessage holds the content sent to each channel.
type NotificationMessage struct {
	// Email fields
	Subject  string
	HtmlBody string
	TextBody string

	// SMS field
	SMSText string

	// Slack field
	SlackText  string
	SlackColor string // "good", "warning", "danger"

	// Push field
	PushTitle string
	PushBody  string
	PushData  map[string]string

	// WebSocket/SSE field
	WSEvent   string
	WSPayload interface{}
}

// Notification is the interface all notification types must implement.
// Inspired by Laravel's Notification class.
type Notification interface {
	// Via returns the list of channel names this notification should be sent through.
	// Valid values: "email", "sms", "slack", "push", "websocket", "sse"
	Via(notifiable Notifiable) []string

	// ToNotification builds the message for the specified channel.
	ToNotification(channel string, notifiable Notifiable) *NotificationMessage
}

// NotificationChannel is the interface each channel driver must implement.
type NotificationChannel interface {
	// Name returns the unique channel identifier (e.g. "email").
	Name() string
	// Send delivers the notification to the notifiable entity.
	Send(ctx context.Context, notifiable Notifiable, notification Notification) error
}

// NotificationManager manages routing notifications to their channels.
type NotificationManager interface {
	// Send dispatches the notification via all channels specified by Via().
	Send(ctx context.Context, notifiable Notifiable, notification Notification) error
	// SendTo sends to a specific subset of channels, overriding Via().
	SendTo(ctx context.Context, notifiable Notifiable, notification Notification, channels ...string) error
}
