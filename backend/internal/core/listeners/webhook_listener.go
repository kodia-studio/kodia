package listeners

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/pkg/webhook"
	"gorm.io/gorm"
)

// WebhookListener listens for any system event and dispatches registered webhooks.
type WebhookListener struct {
	db     *gorm.DB
	engine *webhook.Engine
}

// NewWebhookListener creates a new WebhookListener.
func NewWebhookListener(db *gorm.DB) *WebhookListener {
	return &WebhookListener{
		db:     db,
		engine: webhook.NewEngine(),
	}
}

// Handle processes the event and delivers it to registered webhooks.
func (l *WebhookListener) Handle(event string, data interface{}) error {
	// 1. Fetch registered webhooks for this event
	// In a real app, this would be a table like 'webhooks' (id, url, secret, event)
	// For this upgrade, we'll implement the logic assuming the table exists.
	type Webhook struct {
		ID     string
		URL    string
		Secret string
		Event  string
	}

	webhooks := []Webhook{}
	// Fetch actual webhooks from database. 
	// This will look for a 'webhooks' table in your database.
	if l.db != nil {
		l.db.Where("event = ? OR event = ?", event, "*").Find(&webhooks)
	}

	// Since we are in the Framework Layer, we provide the infrastructure.
	// Let's implement the recording logic in the history table we migrated.
	
	for _, wh := range webhooks {
		go func(w Webhook) {
			code, err := l.engine.Dispatch(w.URL, w.Secret, event, data)
			
			payloadJson, _ := json.Marshal(data)
			success := err == nil && code >= 200 && code < 300
			
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}

			// Record history (Auditing requested by USER)
			l.db.Exec(`
				INSERT INTO webhook_histories (id, url, event, payload, response_code, success, error_message, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			`, fmt.Sprintf("%d", time.Now().UnixNano()), w.URL, event, string(payloadJson), code, success, errorMessage)
			
		}(wh)
	}

	return nil
}

// time import was needed for the Exec, I might have missed it in thinking. 
// I'll add the necessary imports in the actual file.
