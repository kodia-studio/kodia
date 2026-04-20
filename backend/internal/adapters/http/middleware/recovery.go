package middleware

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from panics and logs them.
// Also reports the panic to Sentry if initialized.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log to zap
				log.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Report to Sentry
				hub := sentry.CurrentHub()
				if hub.Client() != nil {
					hub.RecoverWithContext(c.Request.Context(), err)
				}

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
