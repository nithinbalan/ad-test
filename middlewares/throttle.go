package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ThrottleMiddleware creates a middleware that limits request throughput to n requests per second with a burst of m requests.
func ThrottleMiddleware(n int, m int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(n), m)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
