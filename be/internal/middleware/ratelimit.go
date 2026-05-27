package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// RateLimit returns a rate limiting middleware that applies globally.
// It sets rate limit headers on all responses and returns 429 when limit exceeded.
func RateLimit(limiter *helpers.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP (handle proxied requests)
		ip := getClientIP(c)

		// Check rate limit
		info := limiter.Allow(ip)

		// Set rate limit headers on ALL responses
		c.Header("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(info.Reset, 10))

		if !info.Allowed {
			// Calculate retry-after in seconds
			retryAfter := int(info.Reset - time.Now().Unix())
			if retryAfter < 0 {
				retryAfter = 0
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}

// getClientIP extracts the real client IP from the request,
// considering X-Forwarded-For and X-Real-IP headers for proxied requests.
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header (can contain multiple IPs)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		// Take the first IP (original client)
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}
