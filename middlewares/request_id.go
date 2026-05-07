package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the gin.Context key under which the request id is stored.
const RequestIDKey = "request_id"

// requestIDHeader is the inbound/outbound HTTP header name carrying the id.
const requestIDHeader = "X-Request-ID"

// RequestID assigns a unique id to every request, taken from the inbound
// X-Request-ID header when present (so callers can propagate one) or freshly
// generated otherwise. The id is exposed back to the client via the same
// response header and stored on the gin.Context so downstream middlewares
// (api request log, recovery) and handlers can correlate logs with it.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(requestIDHeader)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(RequestIDKey, rid)
		c.Writer.Header().Set(requestIDHeader, rid)
		c.Next()
	}
}

// CurrentRequestID returns the request id that was set on the context by
// RequestID. Returns an empty string if no id is present.
func CurrentRequestID(c *gin.Context) string {
	if v, ok := c.Get(RequestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
