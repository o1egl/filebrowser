package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/filebrowser/filebrowser/v3/log"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

const (
	// RequestIDKey is the key that holds the unique request ID in a request context.
	RequestIDKey ctxKeyRequestID = 0

	RequestIDLogName = "request_id"
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := c.Request.Header.Get(RequestIDHeader)
	if requestID == "" {
		requestID = uuid.New().String()
	}
	ctx = context.WithValue(ctx, RequestIDKey, requestID)
	ctx = log.NewContext(ctx, log.Fields{RequestIDLogName: requestID})
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
