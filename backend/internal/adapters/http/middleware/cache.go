package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kodia-studio/kodia/internal/core/ports"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ResponseCache is a high-performance Redis-backed middleware for GET request caching.
func ResponseCache(cache ports.CacheProvider, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Skip if user is authenticated (unless specifically designed for private cache)
		// For now, let's assume this is public cache for performance.
		if _, exists := c.Get("user"); exists {
			c.Next()
			return
		}

		key := fmt.Sprintf("kodia:cache:%s:%s", c.Request.URL.Path, c.Request.URL.RawQuery)
		
		var cachedResponse struct {
			Body        []byte              `json:"body"`
			Status      int                 `json:"status"`
			Headers     map[string][]string `json:"headers"`
			ContentType string              `json:"content_type"`
		}

		if err := cache.Get(c.Request.Context(), key, &cachedResponse); err == nil {
			// Cache hit
			c.Header("X-Cache", "HIT")
			for k, values := range cachedResponse.Headers {
				for _, v := range values {
					c.Header(k, v)
				}
			}
			c.Data(cachedResponse.Status, cachedResponse.ContentType, cachedResponse.Body)
			c.Abort()
			return
		}

		// Cache miss
		c.Header("X-Cache", "MISS")
		
		// Use a custom response writer to capture the response
		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		// Cache the response if it was successful (200 OK)
		if c.Writer.Status() == http.StatusOK {
			cachedResponse.Body = w.body.Bytes()
			cachedResponse.Status = c.Writer.Status()
			cachedResponse.ContentType = c.Writer.Header().Get("Content-Type")
			cachedResponse.Headers = make(map[string][]string)
			
			// Copy important headers
			for k, v := range c.Writer.Header() {
				if k != "X-Cache" && k != "Set-Cookie" {
					cachedResponse.Headers[k] = v
				}
			}

			_ = cache.Set(c.Request.Context(), key, cachedResponse, ttl)
		}
	}
}

// CacheControl sets the standard HTTP Cache-Control header.
func CacheControl(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
		}
		c.Next()
	}
}

// ETag adds ETag support based on response body content.
func ETag() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			data := w.body.Bytes()
			hash := md5.Sum(data)
			etag := fmt.Sprintf(`W/"%s"`, hex.EncodeToString(hash[:]))
			
			if c.Request.Header.Get("If-None-Match") == etag {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}

			c.Header("ETag", etag)
		}
	}
}
