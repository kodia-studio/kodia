package observability

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/kodia-studio/kodia/pkg/health"
	"go.uber.org/zap"
)

// PulseMessage represents the data sent over the WebSocket.
type PulseMessage struct {
	Type      string      `json:"type"` // "stats" or "log"
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// LogData represents a filtered log entry for Pulse.
type LogData struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Module  string `json:"module,omitempty"`
}

// PulseManager orchestrates real-time telemetry broadcasting.
type PulseManager struct {
	log         *zap.Logger
	clients     map[chan []byte]bool
	register    chan chan []byte
	unregister  chan chan []byte
	broadcast   chan []byte
	logs        chan LogData
	mu          sync.Mutex
	stopChannel chan struct{}
}

// NewPulseManager creates a new PulseManager instance.
func NewPulseManager(log *zap.Logger) *PulseManager {
	return &PulseManager{
		log:         log,
		clients:     make(map[chan []byte]bool),
		register:    make(chan chan []byte),
		unregister:  make(chan chan []byte),
		broadcast:   make(chan []byte),
		logs:        make(chan LogData, 100),
		stopChannel: make(chan struct{}),
	}
}

// Run starts the Pulse broadcasting loop.
func (pm *PulseManager) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	pm.log.Info("Pulse Manager is running")

	for {
		select {
		case <-ticker.C:
			// Gather stats and broadcast
			stats, err := health.Gather(ctx)
			if err == nil {
				pm.send("stats", stats)
			}

		case client := <-pm.register:
			pm.mu.Lock()
			pm.clients[client] = true
			pm.mu.Unlock()

		case client := <-pm.unregister:
			pm.mu.Lock()
			if _, ok := pm.clients[client]; ok {
				delete(pm.clients, client)
				close(client)
			}
			pm.mu.Unlock()

		case message := <-pm.broadcast:
			pm.mu.Lock()
			for client := range pm.clients {
				select {
				case client <- message:
				default:
					close(client)
					delete(pm.clients, client)
				}
			}
			pm.mu.Unlock()

		case logEntry := <-pm.logs:
			pm.send("log", logEntry)

		case <-pm.stopChannel:
			return
		case <-ctx.Done():
			return
		}
	}
}

// Register adds a new client channel.
func (pm *PulseManager) Register() chan []byte {
	ch := make(chan []byte, 256)
	pm.register <- ch
	return ch
}

// Unregister removes a client channel.
func (pm *PulseManager) Unregister(ch chan []byte) {
	pm.unregister <- ch
}

// Log appends a new log entry to be streamed.
func (pm *PulseManager) Log(level, message string) {
	select {
	case pm.logs <- LogData{Level: level, Message: message}:
	default:
		// Drop if buffer full
	}
}

func (pm *PulseManager) send(msgType string, data interface{}) {
	msg := PulseMessage{
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	pm.broadcast <- payload
}

// Stop shuts down the manager.
func (pm *PulseManager) Stop() {
	close(pm.stopChannel)
}
