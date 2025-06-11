package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Timestamp    string `json:"timestamp"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Status       int    `json:"status"`
	DurationMS   int64  `json:"duration_ms"`
	AccountID    string `json:"account_id,omitempty"`
	Email        string `json:"email,omitempty"`
	IP           string `json:"ip"`
	UserAgent    string `json:"user_agent"`
	RequestSize  int    `json:"request_size"`
	ResponseSize int    `json:"response_size"`
}

// Custom response writer to capture response size
type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time
		startTime := time.Now()

		// Get request size
		var requestSize int
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestSize = len(bodyBytes)
			// Restore the body for the next handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create custom response writer to track response size
		customWriter := &responseWriter{ResponseWriter: c.Writer, size: 0}
		c.Writer = customWriter

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Get user info from JWT context (if available)
		accountID, _ := c.Get("account_id")
		email, _ := c.Get("email")

		// Create log entry
		logEntry := LogEntry{
			Timestamp:    startTime.Format(time.RFC3339),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Status:       c.Writer.Status(),
			DurationMS:   duration.Milliseconds(),
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestSize:  requestSize,
			ResponseSize: customWriter.size,
		}

		// Add account info if available (from JWT)
		if accountID != nil {
			logEntry.AccountID = accountID.(string)
		}
		if email != nil {
			logEntry.Email = email.(string)
		}

		// Convert to JSON and log
		logJSON, _ := json.Marshal(logEntry)
		log.Printf("API_REQUEST: %s", string(logJSON))
	}
}
