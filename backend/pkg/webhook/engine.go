package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

// Payload represents the webhook data structure.
type Payload struct {
	Event     string      `json:"event"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Engine handles the signing and dispatching of webhooks.
type Engine struct {
	httpClient *http.Client
}

// NewEngine creates a new Webhook Engine.
func NewEngine() *Engine {
	return &Engine{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Sign generates an HMAC-SHA256 signature for the payload.
func (e *Engine) Sign(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

// Dispatch sends the webhook to the specified URL.
func (e *Engine) Dispatch(url string, secret string, event string, data interface{}) (int, error) {
	payload := Payload{
		Event:     event,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Kodia-Signature", e.Sign(body, secret))
	req.Header.Set("User-Agent", "Kodia-Webhook-Engine/1.0")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
