package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/pkg/response"
)

// IPWhitelist returns a middleware that restricts access to a list of allowed IP addresses or CIDR ranges.
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		isAllowed := false
		for _, allowed := range allowedIPs {
			// Check if it's a CIDR range
			if strings.Contains(allowed, "/") {
				_, ipNet, err := net.ParseCIDR(allowed)
				if err == nil {
					if ipNet.Contains(net.ParseIP(clientIP)) {
						isAllowed = true
						break
					}
				}
			} else if allowed == clientIP {
				isAllowed = true
				break
			}
		}

		if !isAllowed && len(allowedIPs) > 0 {
			response.Forbidden(c, "Access denied: IP address not whitelisted")
			c.Abort()
			return
		}

		c.Next()
	}
}
