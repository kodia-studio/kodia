package ports

import (
	"context"
	"time"
)

// CacheProvider defines the interface for caching operations.
type CacheProvider interface {
	// Get retrieves a value from the cache.
	Get(ctx context.Context, key string, dest interface{}) error
	// Set stores a value in the cache with a TTL.
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	// Delete removes a value from the cache.
	Delete(ctx context.Context, key string) error
	// Remember attempts to get from cache, or runs fn and stores result if miss.
	Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error
}
