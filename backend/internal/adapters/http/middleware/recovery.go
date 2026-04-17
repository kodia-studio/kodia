package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from panics and logs them.
// Returns 500 Internal Server Error instead of crashing the server.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "An unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
