package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"github.com/hibiken/asynq"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"go.uber.org/zap"
)

// InternalDispatcher handles event routing to listeners.
type InternalDispatcher struct {
	listeners map[string][]ports.Listener
	mu        sync.RWMutex
	queue     ports.QueueProvider
	log       *zap.Logger
}

// NewDispatcher creates a new InternalDispatcher.
func NewDispatcher(queue ports.QueueProvider, log *zap.Logger) *InternalDispatcher {
	return &InternalDispatcher{
		listeners: make(map[string][]ports.Listener),
		queue:     queue,
		log:       log,
	}
}

func (d *InternalDispatcher) Register(eventName string, listeners ...ports.Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], listeners...)
}

func (d *InternalDispatcher) Dispatch(ctx context.Context, event ports.Event) error {
	d.mu.RLock()
	listeners, ok := d.listeners[event.Name()]
	d.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, listener := range listeners {
		// Check if it's an async listener
		if al, ok := listener.(ports.AsyncListener); ok && al.ShouldQueue() {
			if err := d.dispatchToQueue(ctx, event, al); err != nil {
				d.log.Error("Failed to dispatch event to queue", 
					zap.String("event", event.Name()), 
					zap.Error(err),
				)
			}
			continue
		}

		// Run sync
		if err := listener.Handle(ctx, event); err != nil {
			d.log.Error("Listener error", 
				zap.String("event", event.Name()), 
				zap.Error(err),
			)
		}
	}

	return nil
}

func (d *InternalDispatcher) dispatchToQueue(ctx context.Context, event ports.Event, listener ports.AsyncListener) error {
	payload, err := json.Marshal(map[string]interface{}{
		"event_name": event.Name(),
		"payload":    event.Payload(),
		"listener":   fmt.Sprintf("%T", listener),
	})
	if err != nil {
		return err
	}

	task := ports.Task{
		Type:    "event.listener.job",
		Payload: payload,
	}

	return d.queue.Enqueue(ctx, task)
}

// HandleListenerTask is the worker handler for async event listeners.
func (d *InternalDispatcher) HandleListenerTask(ctx context.Context, task *asynq.Task) error {
	var data struct {
		EventName string      `json:"event_name"`
		Payload   interface{} `json:"payload"`
		Listener  string      `json:"listener"`
	}

	if err := json.Unmarshal(task.Payload(), &data); err != nil {
		return err
	}

	d.mu.RLock()
	listeners, ok := d.listeners[data.EventName]
	d.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, listener := range listeners {
		if fmt.Sprintf("%T", listener) == data.Listener {
			// Found the listener, create a temporary event for it
			event := &genericEvent{
				name:    data.EventName,
				payload: data.Payload,
			}
			return listener.Handle(ctx, event)
		}
	}

	return fmt.Errorf("listener %s not found for event %s", data.Listener, data.EventName)
}

type genericEvent struct {
	name    string
	payload interface{}
}

func (e *genericEvent) Name() string           { return e.name }
func (e *genericEvent) Payload() interface{}    { return e.payload }
