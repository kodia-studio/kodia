// Package version provides API versioning middleware and helpers for Kodia Framework.
// Supports URL-path versioning (/api/v1/, /api/v2/) with header and Accept header fallback.
package version

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextKey     = "api_version"
	HeaderKey      = "API-Version"
	DefaultVersion = "v1"
)

var (
	pathVersionRe   = regexp.MustCompile(`^v\d+$`)
	acceptVersionRe = regexp.MustCompile(`application/vnd\.kodia\.(v\d+)\+json`)
)

// Middleware extracts the API version and injects it into the Gin context.
// Priority: URL path segment (e.g. /api/v2/users) > API-Version header > Accept header > "v1"
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKey, extract(c))
		c.Next()
	}
}

func extract(c *gin.Context) string {
	// 1. URL path: /api/v2/users → "v2"
	for _, seg := range strings.Split(c.FullPath(), "/") {
		if pathVersionRe.MatchString(seg) {
			return seg
		}
	}
	// 2. API-Version: v2 header
	if h := c.GetHeader(HeaderKey); h != "" {
		return h
	}
	// 3. Accept: application/vnd.kodia.v2+json
	if accept := c.GetHeader("Accept"); accept != "" {
		if m := acceptVersionRe.FindStringSubmatch(accept); len(m) > 1 {
			return m[1]
		}
	}
	return DefaultVersion
}

// Get returns the API version for the current request.
func Get(c *gin.Context) string {
	v, _ := c.Get(ContextKey)
	s, _ := v.(string)
	if s == "" {
		return DefaultVersion
	}
	return s
}

// Since returns true if the current request version >= v (e.g. v2 >= v1).
func Since(c *gin.Context, v string) bool {
	return versionNum(Get(c)) >= versionNum(v)
}

// Until returns true if the current request version <= v.
func Until(c *gin.Context, v string) bool {
	return versionNum(Get(c)) <= versionNum(v)
}

// Deprecate returns a middleware that attaches RFC 8594 deprecation headers.
//
//	sunsetDate: "2026-12-31"
//	alternative: "/api/v2/users"
func Deprecate(sunsetDate, alternative string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Deprecated", "true")
		if sunsetDate != "" {
			if t, err := time.Parse("2006-01-02", sunsetDate); err == nil {
				c.Header("Sunset", t.UTC().Format(http.TimeFormat))
			}
		}
		if alternative != "" {
			c.Header("Link", fmt.Sprintf(`<%s>; rel="successor-version"`, alternative))
		}
		c.Next()
	}
}

func versionNum(v string) int {
	v = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(v), "v"))
	n, _ := strconv.Atoi(v)
	return n
}
