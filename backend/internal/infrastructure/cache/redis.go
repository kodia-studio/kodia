// Package cache provides Redis client implementation for Kodia Framework Port.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kodia-studio/kodia/internal/core/ports"
	"github.com/kodia-studio/kodia/pkg/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisProvider implements ports.CacheProvider using Redis.
type RedisProvider struct {
	client *redis.Client
	log    *zap.Logger
}

// New creates and returns a connected Redis client as ports.CacheProvider.
func New(cfg *config.Config, log *zap.Logger) (ports.CacheProvider, error) {
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
	
	return &RedisProvider{
		client: client,
		log:    log,
	}, nil
}

func (r *RedisProvider) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("cache key %s not found", key)
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisProvider) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisProvider) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisProvider) Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	// 1. Try to get from cache
	err := r.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	// 2. Cache miss, run the function
	val, err := fn()
	if err != nil {
		return err
	}

	// 3. Store in cache
	if err := r.Set(ctx, key, val, ttl); err != nil {
		r.log.Error("Failed to store in cache during Remember", zap.String("key", key), zap.Error(err))
	}

	// 4. Marshal/Unmarshal to ensure dest is populated correctly with the new value
	// (Since Go doesn't have easy generic assignment to interface dest if we want to return the actual type)
	data, _ := json.Marshal(val)
	return json.Unmarshal(data, dest)
}

// GetClient returns the underlying Redis client for advanced operations.
func (r *RedisProvider) GetClient() *redis.Client {
	return r.client
}
