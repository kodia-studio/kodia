// Package cache provides Redis client initialization for Kodia Framework.
package cache

import (
	"context"
	"fmt"

	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// New creates and returns a connected Redis client.
func New(cfg *config.Config, log *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis connected", zap.String("addr", cfg.Redis.Addr()))
	return client, nil
}
