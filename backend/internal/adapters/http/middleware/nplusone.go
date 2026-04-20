package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/performance"
)

// NPlusOne initializes the query tracking context for N+1 detection.
func NPlusOne(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize the tracking context
		ctx := performance.InitContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		// If in debug mode, add a header if N+1 warnings were generated
		if isDebug {
			stats := performance.GetQueryStats(c.Request.Context())
			if stats != nil {
				for table, count := range stats.TableCounts {
					if count >= performance.DefaultNPlusOneThreshold {
						c.Header("X-Kodia-NPlusOne-Warning", fmt.Sprintf("Table '%s' had %d queries", table, count))
					}
				}
				c.Header("X-Kodia-Query-Count", fmt.Sprintf("%d", stats.TotalQueries))
			}
		}
	}
}
