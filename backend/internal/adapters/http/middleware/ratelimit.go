package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimiter implements token bucket rate limiting using Redis.
// Limits requests per IP address with configurable requests and time window.
type RateLimiter struct {
	client      *redis.Client
	maxRequests int64
	windowSecs  int64
	log         *zap.Logger
}

// NewRateLimiter creates a new RateLimiter instance.
func NewRateLimiter(client *redis.Client, maxRequests int64, windowSecs int64, log *zap.Logger) *RateLimiter {
	return &RateLimiter{
		client:      client,
		maxRequests: maxRequests,
		windowSecs:  windowSecs,
		log:         log,
	}
}

// Middleware returns a Gin middleware function for rate limiting.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Identify the requester: User ID (if authenticated) or IP address
		key := getClientIP(c)
		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("user:%v", userID)
		}

		// Check rate limit
		allowed, remaining, retryAfter := rl.checkRateLimit(c.Request.Context(), key)

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(rl.maxRequests, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		if retryAfter > 0 {
			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
		}

		if !allowed {
			rl.log.Warn(
				"Rate limit exceeded",
				zap.String("key", key),
				zap.Int64("limit", rl.maxRequests),
				zap.Int64("window_secs", rl.windowSecs),
				zap.Int64("retry_after", retryAfter),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": fmt.Sprintf("Maximum %d requests per %d seconds", rl.maxRequests, rl.windowSecs),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if request is within rate limit.
// Returns: (allowed, remaining, retryAfter)
func (rl *RateLimiter) checkRateLimit(ctx context.Context, ip string) (bool, int64, int64) {
	key := fmt.Sprintf("ratelimit:%s", ip)
	now := time.Now().Unix()
	windowStart := now - rl.windowSecs

	// Use Lua script for atomic operation (increment and check)
	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window_start = tonumber(ARGV[2])
		local max_requests = tonumber(ARGV[3])
		local window_secs = tonumber(ARGV[4])

		-- Remove old entries outside the window
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)

		-- Count requests in current window
		local current = redis.call('ZCARD', key)

		if current < max_requests then
			-- Add this request
			redis.call('ZADD', key, now, now)
			-- Set expiration
			redis.call('EXPIRE', key, window_secs + 1)
			-- Return allowed with remaining requests
			return {1, max_requests - current - 1, 0}
		else
			-- Rate limit exceeded, calculate retry after
			local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
			local retry_after = 0
			if oldest and oldest[2] then
				retry_after = tonumber(oldest[2]) + window_secs - now
				if retry_after < 1 then
					retry_after = 1
				end
			end
			return {0, 0, retry_after}
		end
	`)

	result, err := script.Run(ctx, rl.client, []string{key}, now, windowStart, rl.maxRequests, rl.windowSecs).Result()
	if err != nil {
		rl.log.Error("Rate limit check failed", zap.Error(err), zap.String("ip", ip))
		// If Redis fails, allow request (fail open)
		return true, rl.maxRequests, 0
	}

	if resultSlice, ok := result.([]interface{}); ok && len(resultSlice) >= 3 {
		allowed := resultSlice[0].(int64) == 1
		remaining := resultSlice[1].(int64)
		retryAfter := resultSlice[2].(int64)
		return allowed, remaining, retryAfter
	}

	// Default to allow if result parsing fails
	return true, rl.maxRequests, 0
}

// getClientIP extracts the client IP address from the request.
// Checks X-Forwarded-For header first (for proxied requests), then falls back to RemoteAddr.
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (set by proxies)
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// Take the first IP from the comma-separated list
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to direct connection IP
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}

// AuthEndpointRateLimiter creates a rate limiter specifically for auth endpoints.
// Limits to 5 requests per 15 minutes per IP.
func AuthEndpointRateLimiter(client *redis.Client, log *zap.Logger) *RateLimiter {
	const (
		maxAuthRequests = 5           // 5 requests
		authWindowSecs  = 15 * 60     // 15 minutes
	)
	return NewRateLimiter(client, maxAuthRequests, authWindowSecs, log)
}

// LooseRateLimiter creates a less strict rate limiter for general endpoints.
// Limits to 100 requests per minute per IP.
func LooseRateLimiter(client *redis.Client, log *zap.Logger) *RateLimiter {
	const (
		maxGeneralRequests = 100      // 100 requests
		generalWindowSecs  = 60       // 1 minute
	)
	return NewRateLimiter(client, maxGeneralRequests, generalWindowSecs, log)
}
