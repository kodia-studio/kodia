package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// TestRateLimiterAllowsRequestsWithinLimit verifies requests are allowed within rate limit
func TestRateLimiterAllowsRequestsWithinLimit(t *testing.T) {
	// Create a test Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Skip test if Redis is not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create rate limiter: 5 requests per 60 seconds
	limiter := NewRateLimiter(client, 5, 60, logger)

	// Clean up any existing test keys
	client.Del(context.Background(), "ratelimit:127.0.0.1")

	// Create a test request
	router := gin.New()
	router.GET("/test", limiter.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Make 5 requests - all should succeed
	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed with status %d, expected 200", i, w.Code)
		}

		remaining := w.Header().Get("X-RateLimit-Remaining")
		t.Logf("Request %d completed, remaining: %s", i, remaining)
	}
}

// TestRateLimiterBlocksRequestsExceedingLimit verifies requests are blocked when limit exceeded
func TestRateLimiterBlocksRequestsExceedingLimit(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create rate limiter: 3 requests per 60 seconds
	limiter := NewRateLimiter(client, 3, 60, logger)

	// Clean up
	client.Del(context.Background(), "ratelimit:127.0.0.1")

	router := gin.New()
	router.GET("/test", limiter.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Make 3 requests - should succeed
	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d failed, expected 200 but got %d", i, w.Code)
		}
	}

	// 4th request should be blocked
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 Too Many Requests, got %d", w.Code)
	}

	// Check retry-after header
	retryAfter := w.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Error("Expected Retry-After header to be set")
	}
}

// TestRateLimiterIsolatesByIP verifies different IPs have separate rate limits
func TestRateLimiterIsolatesByIP(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create rate limiter: 2 requests per 60 seconds
	limiter := NewRateLimiter(client, 2, 60, logger)

	// Clean up
	client.Del(context.Background(), "ratelimit:192.168.1.1", "ratelimit:192.168.1.2")

	router := gin.New()
	router.GET("/test", limiter.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	// Make 2 requests from IP1 - should succeed
	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("IP1 request %d failed", i)
		}
	}

	// 3rd request from IP1 should be blocked
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:8080"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Error("IP1 third request should be blocked")
	}

	// But IP2 should still be able to make 2 requests
	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("IP2 request %d failed", i)
		}
	}
}

// TestGetClientIPExtractsIPCorrectly verifies IP extraction from various sources
func TestGetClientIPExtractsIPCorrectly(t *testing.T) {
	tests := []struct {
		name         string
		remoteAddr   string
		xForwardedFor string
		xRealIP      string
		expectedIP   string
	}{
		{
			name:       "direct connection",
			remoteAddr: "192.168.1.100:8080",
			expectedIP: "192.168.1.100",
		},
		{
			name:         "x-forwarded-for single IP",
			xForwardedFor: "203.0.113.195",
			remoteAddr:    "192.168.1.100:8080",
			expectedIP:    "203.0.113.195",
		},
		{
			name:         "x-forwarded-for multiple IPs",
			xForwardedFor: "203.0.113.195, 70.41.3.18, 150.172.238.178",
			remoteAddr:    "192.168.1.100:8080",
			expectedIP:    "203.0.113.195",
		},
		{
			name:       "x-real-ip",
			xRealIP:    "203.0.113.195",
			remoteAddr: "192.168.1.100:8080",
			expectedIP: "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.RemoteAddr = tt.remoteAddr

			if tt.xForwardedFor != "" {
				c.Request.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				c.Request.Header.Set("X-Real-IP", tt.xRealIP)
			}

			ip := getClientIP(c)
			if ip != tt.expectedIP {
				t.Errorf("Got IP %s, expected %s", ip, tt.expectedIP)
			}
		})
	}
}

// TestAuthEndpointRateLimiter verifies auth-specific rate limiter configuration
func TestAuthEndpointRateLimiter(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create auth rate limiter (5 requests per 15 minutes)
	limiter := AuthEndpointRateLimiter(client, logger)

	// Verify configuration
	if limiter.maxRequests != 5 {
		t.Errorf("Expected maxRequests=5, got %d", limiter.maxRequests)
	}

	if limiter.windowSecs != 15*60 {
		t.Errorf("Expected windowSecs=900, got %d", limiter.windowSecs)
	}
}

// TestLooseRateLimiter verifies general-purpose rate limiter configuration
func TestLooseRateLimiter(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create loose rate limiter (100 requests per minute)
	limiter := LooseRateLimiter(client, logger)

	// Verify configuration
	if limiter.maxRequests != 100 {
		t.Errorf("Expected maxRequests=100, got %d", limiter.maxRequests)
	}

	if limiter.windowSecs != 60 {
		t.Errorf("Expected windowSecs=60, got %d", limiter.windowSecs)
	}
}

// TestRateLimiterHeadersSet verifies rate limit headers are set correctly
func TestRateLimiterHeadersSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available for testing")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	limiter := NewRateLimiter(client, 3, 60, logger)
	client.Del(context.Background(), "ratelimit:127.0.0.1")

	router := gin.New()
	router.GET("/test", limiter.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check headers
	if w.Header().Get("X-RateLimit-Limit") != "3" {
		t.Error("X-RateLimit-Limit header not set correctly")
	}

	if w.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("X-RateLimit-Remaining header not set")
	}
}

// BenchmarkRateLimiter measures rate limiting performance
func BenchmarkRateLimiter(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		b.Skip("Redis not available for benchmarking")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	limiter := NewRateLimiter(client, 1000, 60, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		allowed, _, _ := limiter.checkRateLimit(context.Background(), fmt.Sprintf("192.168.1.%d", i%254+1))
		if !allowed && i < 1000 {
			b.Fatalf("Rate limit check failed unexpectedly")
		}
	}
}
