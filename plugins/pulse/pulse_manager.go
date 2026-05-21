package pulse

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/kodia-studio/kodia/pkg/health"
	"go.uber.org/zap"
)

type Message struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type LogData struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Module  string `json:"module,omitempty"`
}

type Manager struct {
	log         *zap.Logger
	clients     map[chan []byte]bool
	register    chan chan []byte
	unregister  chan chan []byte
	broadcast   chan []byte
	logs        chan LogData
	mu          sync.Mutex
	stopChannel chan struct{}
}

func NewManager(log *zap.Logger) *Manager {
	return &Manager{
		log:         log,
		clients:     make(map[chan []byte]bool),
		register:    make(chan chan []byte),
		unregister:  make(chan chan []byte),
		broadcast:   make(chan []byte),
		logs:        make(chan LogData, 100),
		stopChannel: make(chan struct{}),
	}
}

func (pm *Manager) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	pm.log.Info("Pulse Manager is running")

	for {
		select {
		case <-ticker.C:
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

func (pm *Manager) Register() chan []byte {
	ch := make(chan []byte, 256)
	pm.register <- ch
	return ch
}

func (pm *Manager) Unregister(ch chan []byte) {
	pm.unregister <- ch
}

func (pm *Manager) Log(level, message string) {
	select {
	case pm.logs <- LogData{Level: level, Message: message}:
	default:
	}
}

func (pm *Manager) send(msgType string, data interface{}) {
	msg := Message{
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

func (pm *Manager) Stop() {
	close(pm.stopChannel)
}
