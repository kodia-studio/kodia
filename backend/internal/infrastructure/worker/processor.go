package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/kodia-studio/kodia/pkg/config"
	"go.uber.org/zap"
)

// Processor handles the registration and execution of background jobs.
type Processor struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	log    *zap.Logger
}

// NewProcessor creates a new asynq Server and Mux.
func NewProcessor(cfg *config.Config, log *zap.Logger) *Processor {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: 10,
			Logger:      &asynqLogger{log: log},
		},
	)

	return &Processor{
		server: srv,
		mux:    asynq.NewServeMux(),
		log:    log,
	}
}

// Start starts the worker server.
func (p *Processor) Start() error {
	p.log.Info("Starting worker processor...")
	return p.server.Run(p.mux)
}

// Register adds a handler for a specific task type.
func (p *Processor) Register(taskType string, handler func(context.Context, *asynq.Task) error) {
	p.mux.HandleFunc(taskType, handler)
}

// asynqLogger adapts zap.Logger to asynq.Logger interface.
type asynqLogger struct {
	log *zap.Logger
}

func (l *asynqLogger) Debug(args ...interface{}) { l.log.Sugar().Debug(args...) }
func (l *asynqLogger) Info(args ...interface{})  { l.log.Sugar().Info(args...) }
func (l *asynqLogger) Warn(args ...interface{})  { l.log.Sugar().Warn(args...) }
func (l *asynqLogger) Error(args ...interface{}) { l.log.Sugar().Error(args...) }
func (l *asynqLogger) Fatal(args ...interface{}) { l.log.Sugar().Fatal(args...) }
