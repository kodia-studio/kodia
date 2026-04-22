package middleware

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SessionStore handles session persistence with Redis + DB fallback.
type SessionStore struct {
	redis *redis.Client
	db    *gorm.DB
}

// NewSessionStore creates a new SessionStore.
func NewSessionStore(redis *redis.Client, db *gorm.DB) *SessionStore {
	return &SessionStore{redis: redis, db: db}
}

// Get retrieves session data from Redis, falling back to DB.
func (s *SessionStore) Get(ctx context.Context, sessionID string) (map[string]interface{}, error) {
	// 1. Try Redis
	if s.redis != nil {
		val, err := s.redis.Get(ctx, "session:"+sessionID).Result()
		if err == nil {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(val), &data); err == nil {
				return data, nil
			}
		}
	}

	// 2. Try DB fallback
	// This would query the 'sessions' table created earlier
	return nil, nil // Implementation placeholder
}

// SessionMiddleware provides traditional cookie-based auth.
func SessionMiddleware(store *SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("kodia_session")
		if err != nil {
			c.Next()
			return
		}

		data, err := store.Get(c.Request.Context(), sessionID)
		if err != nil || data == nil {
			c.Next()
			return
		}

		// Inject user_id from session into context
		if userID, ok := data["user_id"].(string); ok {
			c.Set("user_id", userID)
			c.Set("auth_method", "session")
		}

		c.Next()
	}
}

// RequireAuth ensures the request is authenticated via any method.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user_id"); !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}
		c.Next()
	}
}
