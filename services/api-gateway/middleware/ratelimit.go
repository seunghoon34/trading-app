package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int           // requests per window
	window   time.Duration // time window
}

type RateLimitLogEntry struct {
	Timestamp    string `json:"timestamp"`
	Event        string `json:"event"`
	Identifier   string `json:"identifier"`
	RequestCount int    `json:"request_count"`
	Limit        int    `json:"limit"`
	Path         string `json:"path"`
	LimitType    string `json:"limit_type"` // e.g., "general", "trading", "auth"
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine to prevent memory leaks
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			// Remove old requests outside the window
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < rl.window {
					validTimes = append(validTimes, t)
				}
			}
			if len(validTimes) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validTimes
			}
		}
		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) Allow(identifier string) (bool, int) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Get current requests for this identifier
	times := rl.requests[identifier]

	// Remove old requests outside the window
	var validTimes []time.Time
	for _, t := range times {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	// Check if under limit
	if len(validTimes) >= rl.limit {
		rl.requests[identifier] = validTimes
		return false, len(validTimes)
	}

	// Add current request
	validTimes = append(validTimes, now)
	rl.requests[identifier] = validTimes

	return true, len(validTimes)
}

// Global rate limiters for different endpoints
var (
	// General API rate limiter: 100 requests per minute
	generalLimiter = NewRateLimiter(100, time.Minute)

	// Trading endpoint limiter: 30 requests per minute (more restrictive)
	tradingLimiter = NewRateLimiter(30, time.Minute)

	// Auth endpoint limiter: 10 requests per minute (prevent brute force)
	authLimiter = NewRateLimiter(10, time.Minute)
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine identifier (prefer account_id from JWT, fallback to IP)
		identifier := c.ClientIP()
		if accountID, exists := c.Get("account_id"); exists {
			identifier = "user:" + accountID.(string)
		} else {
			identifier = "ip:" + identifier
		}

		// Choose appropriate rate limiter based on path
		var limiter *RateLimiter
		var limitType string

		path := c.Request.URL.Path
		switch {
		case strings.Contains(path, "/api/v1/auth/"):
			limiter = authLimiter
			limitType = "auth"
		case strings.Contains(path, "/api/v1/trading/"):
			limiter = tradingLimiter
			limitType = "trading"
		default:
			limiter = generalLimiter
			limitType = "general"
		}

		// Check rate limit
		allowed, currentCount := limiter.Allow(identifier)

		if !allowed {
			// Log rate limit violation
			logEntry := RateLimitLogEntry{
				Timestamp:    time.Now().Format(time.RFC3339),
				Event:        "rate_limit_exceeded",
				Identifier:   identifier,
				RequestCount: currentCount,
				Limit:        limiter.limit,
				Path:         path,
				LimitType:    limitType,
			}
			logJSON, _ := json.Marshal(logEntry)
			log.Printf("RATE_LIMIT: %s", string(logJSON))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"limit":       limiter.limit,
				"window":      limiter.window.String(),
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", string(rune(limiter.limit)))
		c.Header("X-RateLimit-Remaining", string(rune(limiter.limit-currentCount)))
		c.Header("X-RateLimit-Window", limiter.window.String())

		c.Next()
	}
}
