package middleware

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
	"github.com/richardlinio/waterballsa-backend-golang/internal/config"
)

// ipRecord stores request count and window start time for an IP address
type ipRecord struct {
	count       int
	windowStart time.Time
	mu          sync.RWMutex
}

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	requestsPerMinute int
	windowDuration    time.Duration
	records           sync.Map // map[string]*ipRecord
	cleanupInterval   time.Duration
	recordExpiry      time.Duration
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config config.RateLimitConfig) *RateLimiter {
	limiter := &RateLimiter{
		requestsPerMinute: config.RequestsPerMinute,
		windowDuration:    config.WindowDuration,
		cleanupInterval:   config.CleanupInterval,
		recordExpiry:      config.RecordExpiry,
	}

	// Start background cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()

	// Load or create record for this IP
	value, _ := rl.records.LoadOrStore(ip, &ipRecord{
		count:       0,
		windowStart: now,
	})

	record, ok := value.(*ipRecord)
	if !ok {
		// This should never happen, but handle it gracefully
		return false
	}
	record.mu.Lock()
	defer record.mu.Unlock()

	// Check if we need to reset the window
	if now.Sub(record.windowStart) >= rl.windowDuration {
		record.count = 0
		record.windowStart = now
	}

	// Check if request is allowed
	if record.count >= rl.requestsPerMinute {
		return false
	}

	// Increment count and allow request
	record.count++
	return true
}

// cleanup periodically removes expired IP records to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		rl.records.Range(func(key, value any) bool {
			record, ok := value.(*ipRecord)
			if !ok {
				// Invalid record type, delete it
				rl.records.Delete(key)
				return true
			}

			record.mu.RLock()
			windowStart := record.windowStart
			record.mu.RUnlock()

			// Remove records that haven't been used for more than recordExpiry
			if now.Sub(windowStart) > rl.recordExpiry {
				rl.records.Delete(key)
			}
			return true
		})
	}
}

// RateLimit returns a Gin middleware that enforces rate limiting
func RateLimit(config config.RateLimitConfig, logger *slog.Logger) gin.HandlerFunc {
	// Create rate limiter instance
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		// Skip if rate limiting is disabled
		if !config.Enabled {
			c.Next()
			return
		}

		// Skip health check endpoint
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		// Get client IP
		ip := c.ClientIP()

		// Check if request is allowed
		if !limiter.Allow(ip) {
			_ = c.Error(apperror.RateLimitExceeded())
			c.Abort()
			return
		}

		c.Next()
	}
}
