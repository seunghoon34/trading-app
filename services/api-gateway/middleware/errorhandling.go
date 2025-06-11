package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorLogEntry struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Error     string `json:"error"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	IP        string `json:"ip"`
	Stack     string `json:"stack,omitempty"`
}

type StandardErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				logEntry := ErrorLogEntry{
					Timestamp: time.Now().Format(time.RFC3339),
					Event:     "panic_recovered",
					Error:     err.(string),
					Path:      c.Request.URL.Path,
					Method:    c.Request.Method,
					IP:        c.ClientIP(),
					Stack:     string(debug.Stack()),
				}

				logJSON, _ := json.Marshal(logEntry)
				log.Printf("ERROR: %s", string(logJSON))

				// Return standardized error response
				errorResponse := StandardErrorResponse{
					Error:     "internal_server_error",
					Message:   "An internal server error occurred",
					Timestamp: time.Now().Format(time.RFC3339),
					Path:      c.Request.URL.Path,
				}

				c.JSON(http.StatusInternalServerError, errorResponse)
				c.Abort()
			}
		}()

		// Process request
		c.Next()

		// Check for errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log the error
			logEntry := ErrorLogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Event:     "request_error",
				Error:     err.Error(),
				Path:      c.Request.URL.Path,
				Method:    c.Request.Method,
				IP:        c.ClientIP(),
			}

			logJSON, _ := json.Marshal(logEntry)
			log.Printf("ERROR: %s", string(logJSON))

			// If no response was sent yet, send standardized error
			if !c.Writer.Written() {
				errorResponse := StandardErrorResponse{
					Error:     "request_failed",
					Message:   "The request could not be processed",
					Timestamp: time.Now().Format(time.RFC3339),
					Path:      c.Request.URL.Path,
				}

				c.JSON(http.StatusInternalServerError, errorResponse)
			}
		}
	}
}
