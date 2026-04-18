package ports

import (
	"context"
	"time"
)

// Task represents a background job.
type Task struct {
	Type    string
	Payload []byte
}

// QueueProvider defines the interface for enqueueing background jobs.
type QueueProvider interface {
	// Enqueue adds a task to the default queue.
	Enqueue(ctx context.Context, task Task, opts ...interface{}) error
	// EnqueueAt schedules a task to be executed at a specific time.
	EnqueueAt(ctx context.Context, task Task, at time.Time, opts ...interface{}) error
	// Close closes the queue provider connection.
	Close() error
}
