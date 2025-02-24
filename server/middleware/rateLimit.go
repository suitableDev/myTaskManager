package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var (
	SignupLimiter = ratelimit.NewBucketWithRate(5, 10)
	LoginLimiter  = ratelimit.NewBucketWithRate(10, 20)
	CreateLimiter = ratelimit.NewBucketWithRate(10, 20)
	DeleteLimiter = ratelimit.NewBucketWithRate(10, 20)
	HealthLimiter = ratelimit.NewBucketWithRate(10, 50)
)

func RateLimitMiddleware(rateLimiter *ratelimit.Bucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimiter.TakeAvailable(1) == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}
