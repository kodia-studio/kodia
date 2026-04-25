package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kodia-studio/kodia/internal/providers"
)

// MailQueue handles queuing emails for async sending
type MailQueue struct {
	client   *asynq.Client
	provider providers.MailProvider
}

// MailTask represents a mail task in the queue
type MailTask struct {
	Mail *providers.Mail
}

const (
	// Queue name
	MailQueueName = "mail"
	// Task type
	SendMailTask = "mail:send"
)

// NewMailQueue creates a new mail queue
func NewMailQueue(redisAddr string, provider providers.MailProvider) (*MailQueue, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	return &MailQueue{
		client:   client,
		provider: provider,
	}, nil
}

// Enqueue enqueues a mail message for sending
func (mq *MailQueue) Enqueue(ctx context.Context, mail *providers.Mail) error {
	payload, err := json.Marshal(MailTask{Mail: mail})
	if err != nil {
		return fmt.Errorf("failed to marshal mail task: %w", err)
	}

	task := asynq.NewTask(SendMailTask, payload)

	info, err := mq.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue mail task: %w", err)
	}

	fmt.Printf("Mail enqueued: ID=%s, Queue=%s, MaxRetry=%d\n",
		info.ID, info.Queue, info.MaxRetry)

	return nil
}

// EnqueueWithOptions enqueues a mail message with custom options
func (mq *MailQueue) EnqueueWithOptions(ctx context.Context, mail *providers.Mail, opts ...asynq.Option) error {
	payload, err := json.Marshal(MailTask{Mail: mail})
	if err != nil {
		return fmt.Errorf("failed to marshal mail task: %w", err)
	}

	task := asynq.NewTask(SendMailTask, payload, opts...)

	_, err = mq.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue mail task with options: %w", err)
	}

	return nil
}

// EnqueueDelayed enqueues a mail message to be sent after delay
func (mq *MailQueue) EnqueueDelayed(ctx context.Context, mail *providers.Mail, delay time.Duration) error {
	return mq.EnqueueWithOptions(ctx, mail, asynq.ProcessIn(delay))
}

// Close closes the mail queue
func (mq *MailQueue) Close() error {
	return mq.client.Close()
}

// MailServer handles processing mail tasks from the queue
type MailServer struct {
	mux      *asynq.ServeMux
	provider providers.MailProvider
}

// NewMailServer creates a new mail server for processing queued tasks
func NewMailServer(provider providers.MailProvider) *MailServer {
	mux := asynq.NewServeMux()
	server := &MailServer{
		mux:      mux,
		provider: provider,
	}

	// Register mail task handler
	mux.HandleFunc(SendMailTask, server.handleSendMail)

	return server
}

// handleSendMail handles sending mail tasks
func (ms *MailServer) handleSendMail(ctx context.Context, t *asynq.Task) error {
	var task MailTask

	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", string(t.Payload()), err)
	}

	if err := task.Mail.Validate(); err != nil {
		return fmt.Errorf("mail validation failed: %w", err)
	}

	// Send mail using provider
	if err := ms.provider.Send(ctx, task.Mail); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	fmt.Printf("Mail sent successfully: To=%v\n", task.Mail.To)
	return nil
}

// GetMux returns the asynq serve mux
func (ms *MailServer) GetMux() *asynq.ServeMux {
	return ms.mux
}

// ProcessMailQueue starts processing mail queue
func ProcessMailQueue(redisAddr string, provider providers.MailProvider, concurrency int) error {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				MailQueueName: 10, // Increase concurrency for mail queue
			},
		},
	)

	mux := asynq.NewServeMux()
	mailServer := NewMailServer(provider)

	mux.HandleFunc(SendMailTask, mailServer.handleSendMail)

	if err := server.Run(mux); err != nil {
		return fmt.Errorf("failed to run mail server: %w", err)
	}

	return nil
}
