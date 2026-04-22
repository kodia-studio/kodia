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
	// EnqueueChain enqueues a sequence of tasks to be executed one after another.
	EnqueueChain(ctx context.Context, tasks ...Task) error
	// EnqueueBatch enqueues a group of tasks to be processed as a single unit.
	EnqueueBatch(ctx context.Context, tasks []Task) error
	// Close closes the queue provider connection.
	Close() error
}
