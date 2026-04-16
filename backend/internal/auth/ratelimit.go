package auth

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter provides per-IP rate limiting using token buckets.
type IPRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitorEntry
	rate     rate.Limit
	burst    int
	ttl      time.Duration
}

type visitorEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter creates a new per-IP rate limiter.
// r is tokens per second, burst is the bucket size.
func NewIPRateLimiter(r float64, burst int) *IPRateLimiter {
	rl := &IPRateLimiter{
		visitors: make(map[string]*visitorEntry),
		rate:     rate.Limit(r),
		burst:    burst,
		ttl:      3 * time.Minute,
	}
	go rl.cleanup()
	return rl
}

func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if entry, exists := rl.visitors[ip]; exists {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.visitors[ip] = &visitorEntry{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

func (rl *IPRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, entry := range rl.visitors {
			if time.Since(entry.lastSeen) > rl.ttl {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a Gin middleware that rate-limits by client IP.
func RateLimitMiddleware(r float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests, please try again later"})
			c.Abort()
			return
		}
		c.Next()
	}
}
