package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"

	helper "task-manager/server/helpers"
)

var rateLimiters = make(map[string]map[string]*ratelimit.Bucket)
var mu sync.Mutex

// Get a rate limiter for a specific route and IP
func getRateLimiter(route string, ip string, rate float64, capacity int64) *ratelimit.Bucket {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := rateLimiters[route]; !exists {
		rateLimiters[route] = make(map[string]*ratelimit.Bucket)
	}

	if _, exists := rateLimiters[route][ip]; !exists {
		rateLimiters[route][ip] = ratelimit.NewBucketWithRate(rate, capacity)
	}

	return rateLimiters[route][ip]
}

// RateLimitMiddleware creates a rate limiter for a specific route
func RateLimitMiddleware(rate float64, capacity int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		route := c.FullPath()
		bucket := getRateLimiter(route, clientIP, rate, capacity)

		if bucket.TakeAvailable(1) == 0 {
			helper.RespondWithError(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests for this route")
			c.Abort()
			return
		}

		c.Next()
	}
}
