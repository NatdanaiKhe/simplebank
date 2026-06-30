package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
	requestIDKey    = "requestID"
	loggerKey       = "logger"
)

// RequestID assigns a unique identifier to every request. If the client
// sends an X-Request-ID header it is reused (distributed tracing); otherwise
// a new UUID v4 is generated. The ID is written to the response header and
// stored in the gin context for handlers and loggers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Header(requestIDHeader, id)
		c.Set(requestIDKey, id)
		c.Next()
	}
}

// GetRequestID extracts the request ID from the gin context.
// Returns an empty string if the middleware was not registered.
func GetRequestID(c *gin.Context) string {
	id, exists := c.Get(requestIDKey)
	if !exists {
		return ""
	}
	return id.(string)
}

// LoggerMiddleware injects a request-scoped zap.Logger into the gin context.
// The logger has the request_id field pre-attached so every log line from
// handlers is automatically traceable. Must run after RequestID middleware.
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLogger := logger.With(zap.String("request_id", GetRequestID(c)))
		c.Set(loggerKey, requestLogger)
		c.Next()
	}
}

// GetLogger returns the request-scoped logger from the gin context.
// Falls back to the global logger if the middleware was not registered.
func GetLogger(c *gin.Context) *zap.Logger {
	l, exists := c.Get(loggerKey)
	if !exists {
		return zap.L()
	}
	return l.(*zap.Logger)
}
