package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the context key for request ID
const RequestIDKey = "requestID"

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Get request ID
		requestID, _ := c.Get(RequestIDKey)
		reqID, _ := requestID.(string)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		status := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user ID if authenticated
		userID := ""
		if uid, exists := c.Get(UserIDKey); exists {
			userID, _ = uid.(string)
		}

		// Build log message
		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%s] %s %s %d %s %s user=%s err=%s",
			reqID,
			c.Request.Method,
			path,
			status,
			latency,
			clientIP,
			userID,
			c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}

// bodyLogWriter wraps gin.ResponseWriter to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// DetailedLoggingMiddleware logs requests with body (use only in debug mode)
func DetailedLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Log details
		latency := time.Since(start)

		log.Printf(`
========== REQUEST ==========
Method: %s
Path: %s
Headers: %v
Body: %s
========== RESPONSE ==========
Status: %d
Latency: %s
Body: %s
=============================`,
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Header,
			string(requestBody),
			c.Writer.Status(),
			latency,
			blw.body.String(),
		)
	}
}

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get(RequestIDKey)
				reqID, _ := requestID.(string)

				log.Printf("[%s] PANIC recovered: %v", reqID, err)

				c.AbortWithStatusJSON(500, gin.H{
					"error":      "Internal Server Error",
					"message":    "An unexpected error occurred",
					"request_id": reqID,
				})
			}
		}()
		c.Next()
	}
}

