package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders returns a middleware that sets various security-enhancing HTTP headers.
// These headers help protect against common web vulnerabilities like XSS, Clickjacking, and Sniffing.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strict-Transport-Security (HSTS)
		// Tells the browser to only interact with the server using HTTPS.
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Content-Security-Policy (CSP)
		// Controls which resources the browser is allowed to load.
		// Default: restricted, allow self and specific common CDNs if needed.
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; object-src 'none'; frame-ancestors 'none'; upgrade-insecure-requests")

		// X-Frame-Options
		// Prevents Clickjacking by not allowing the site to be embedded in an iframe.
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		// Prevents the browser from sniffing the MIME type away from the declared Content-Type.
		c.Header("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		// Enables the browser's XSS filter. (Legacy but still useful for older browsers).
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		// Controls how much referrer information is included with requests.
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy
		// Controls which browser features can be used (camera, microphone, geolocation).
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")

		c.Next()
	}
}
