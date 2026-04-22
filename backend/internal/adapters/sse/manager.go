// Package sse provides Server-Sent Events support for Kodia Framework.
package sse

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SSEClient represents a single SSE subscriber connection.
type SSEClient struct {
	ID      string
	UserID  string
	Channel string
	Events  chan *SSEEvent
}

// SSEEvent represents a single server-sent event.
type SSEEvent struct {
	ID    string
	Event string
	Data  interface{}
}

// Format serialises the event to the SSE wire format.
func (e *SSEEvent) Format() string {
	data, _ := json.Marshal(e.Data)
	out := ""
	if e.ID != "" {
		out += fmt.Sprintf("id: %s\n", e.ID)
	}
	if e.Event != "" {
		out += fmt.Sprintf("event: %s\n", e.Event)
	}
	out += fmt.Sprintf("data: %s\n\n", string(data))
	return out
}

// Manager manages all active SSE client connections.
type Manager struct {
	// clients indexed by client ID
	clients map[string]*SSEClient
	// userClients maps userID -> list of client IDs
	userClients map[string][]string
	// channelClients maps channel -> list of client IDs
	channelClients map[string][]string

	mu      sync.RWMutex
	counter int64
}

// NewManager creates a new SSE Manager.
func NewManager() *Manager {
	return &Manager{
		clients:        make(map[string]*SSEClient),
		userClients:    make(map[string][]string),
		channelClients: make(map[string][]string),
	}
}

// Subscribe registers a new SSE client and returns it.
func (m *Manager) Subscribe(userID, channel string) *SSEClient {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter++
	client := &SSEClient{
		ID:      fmt.Sprintf("sse-%d-%d", time.Now().UnixNano(), m.counter),
		UserID:  userID,
		Channel: channel,
		Events:  make(chan *SSEEvent, 32),
	}

	m.clients[client.ID] = client

	if userID != "" {
		m.userClients[userID] = append(m.userClients[userID], client.ID)
	}
	if channel != "" {
		m.channelClients[channel] = append(m.channelClients[channel], client.ID)
	}

	return client
}

// Unsubscribe removes a client from the manager.
func (m *Manager) Unsubscribe(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, ok := m.clients[clientID]
	if !ok {
		return
	}

	delete(m.clients, clientID)

	// Remove from user index
	if client.UserID != "" {
		ids := m.userClients[client.UserID]
		for i, id := range ids {
			if id == clientID {
				m.userClients[client.UserID] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
		if len(m.userClients[client.UserID]) == 0 {
			delete(m.userClients, client.UserID)
		}
	}

	// Remove from channel index
	if client.Channel != "" {
		ids := m.channelClients[client.Channel]
		for i, id := range ids {
			if id == clientID {
				m.channelClients[client.Channel] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
		if len(m.channelClients[client.Channel]) == 0 {
			delete(m.channelClients, client.Channel)
		}
	}

	close(client.Events)
}

// Publish sends an event to all clients subscribed to a channel.
func (m *Manager) Publish(channel, eventName string, data interface{}) error {
	m.mu.RLock()
	ids := append([]string{}, m.channelClients[channel]...)
	m.mu.RUnlock()

	event := &SSEEvent{
		ID:    fmt.Sprintf("%d", time.Now().UnixMicro()),
		Event: eventName,
		Data:  data,
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, id := range ids {
		if c, ok := m.clients[id]; ok {
			select {
			case c.Events <- event:
			default:
				// Drop if buffer full
			}
		}
	}
	return nil
}

// PublishToUser sends an event to all SSE connections of a specific user.
func (m *Manager) PublishToUser(userID, eventName string, data interface{}) error {
	m.mu.RLock()
	ids := append([]string{}, m.userClients[userID]...)
	m.mu.RUnlock()

	event := &SSEEvent{
		ID:    fmt.Sprintf("%d", time.Now().UnixMicro()),
		Event: eventName,
		Data:  data,
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, id := range ids {
		if c, ok := m.clients[id]; ok {
			select {
			case c.Events <- event:
			default:
			}
		}
	}
	return nil
}

// ClientCount returns the number of active SSE connections.
func (m *Manager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}
