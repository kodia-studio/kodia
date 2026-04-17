package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a structured request logging middleware using Zap.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userID := GetUserID(c)

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		}

		if userID != "" {
			fields = append(fields, zap.String("user_id", userID))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		switch {
		case statusCode >= 500:
			log.Error("Server error", fields...)
		case statusCode >= 400:
			log.Warn("Client error", fields...)
		default:
			log.Info("Request", fields...)
		}
	}
}
