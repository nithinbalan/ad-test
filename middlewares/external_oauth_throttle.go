package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ExternalOAuthThrottle struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	maxReqs  int
	window   time.Duration
}

func NewExternalOAuthThrottle(maxReqs int, window time.Duration) *ExternalOAuthThrottle {
	return &ExternalOAuthThrottle{
		requests: make(map[string][]time.Time),
		maxReqs:  maxReqs,
		window:   window,
	}
}

func (t *ExternalOAuthThrottle) Allow(clientIP string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	times := t.requests[clientIP]

	// Remove expired requests for this IP.
	validTimes := make([]time.Time, 0, len(times))
	for _, ts := range times {
		if ts.After(cutoff) {
			validTimes = append(validTimes, ts)
		}
	}

	// Check rate limit.
	if len(validTimes) >= t.maxReqs {
		t.requests[clientIP] = validTimes
		return false
	}

	// Record current request.
	validTimes = append(validTimes, now)
	t.requests[clientIP] = validTimes

	return true
}

var externalOAuthThrottle = NewExternalOAuthThrottle(5, time.Minute)

func ExternalOAuthThrottleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !externalOAuthThrottle.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":  "TOO_MANY_REQUESTS",
				"error": "Too many authentication attempts. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
