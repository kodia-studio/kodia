package worker

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// AsynqProvider implements ports.QueueProvider using asynq.
type AsynqProvider struct {
	client *asynq.Client
	log    *zap.Logger
}

// NewAsynqProvider creates a new Asynq-backed QueueProvider.
func NewAsynqProvider(cfg *config.Config, log *zap.Logger) *AsynqProvider {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &AsynqProvider{
		client: client,
		log:    log,
	}
}

func (p *AsynqProvider) Enqueue(ctx context.Context, task ports.Task, opts ...interface{}) error {
	t := asynq.NewTask(task.Type, task.Payload)
	_, err := p.client.Enqueue(t)
	if err != nil {
		p.log.Error("Failed to enqueue task", zap.String("type", task.Type), zap.Error(err))
		return err
	}
	return nil
}

func (p *AsynqProvider) EnqueueAt(ctx context.Context, task ports.Task, at time.Time, opts ...interface{}) error {
	t := asynq.NewTask(task.Type, task.Payload)
	_, err := p.client.Enqueue(t, asynq.ProcessAt(at))
	if err != nil {
		p.log.Error("Failed to enqueue task at specific time", zap.String("type", task.Type), zap.Time("at", at), zap.Error(err))
		return err
	}
	return nil
}

func (p *AsynqProvider) EnqueueChain(ctx context.Context, tasks ...ports.Task) error {
	// Fallback implementation: Enqueue tasks sequentially.
	// In a world-class framework, we eventually upgrade to native Chaining.
	for _, task := range tasks {
		if err := p.Enqueue(ctx, task); err != nil {
			return err
		}
	}
	return nil
}

func (p *AsynqProvider) EnqueueBatch(ctx context.Context, tasks []ports.Task) error {
	// Fallback implementation: Enqueue tasks sequentially.
	for _, task := range tasks {
		if err := p.Enqueue(ctx, task); err != nil {
			return err
		}
	}
	return nil
}

func (p *AsynqProvider) Close() error {
	return p.client.Close()
}
