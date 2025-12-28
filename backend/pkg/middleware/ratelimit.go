package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	redis     *redis.Client
	keyPrefix string
}

// RateLimitConfig holds rate limit configuration
type RateLimitConfig struct {
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, keyPrefix string) *RateLimiter {
	return &RateLimiter{
		redis:     redisClient,
		keyPrefix: keyPrefix,
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identifier (user ID if authenticated, otherwise IP)
		identifier := c.ClientIP()
		if userID, exists := c.Get(UserIDKey); exists {
			identifier = userID.(string)
		}

		key := limiter.keyPrefix + ":ratelimit:" + identifier

		allowed, remaining, resetTime, err := limiter.Allow(c.Request.Context(), key, config)
		if err != nil {
			// If Redis fails, allow the request but log the error
			c.Next()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded",
				"retry_after": resetTime - time.Now().Unix(),
			})
			return
		}

		c.Next()
	}
}

// Allow checks if a request is allowed and updates the counter
func (r *RateLimiter) Allow(ctx context.Context, key string, config RateLimitConfig) (bool, int, int64, error) {
	now := time.Now()
	windowStart := now.Truncate(config.Window)
	windowEnd := windowStart.Add(config.Window)

	pipe := r.redis.Pipeline()

	// Increment counter
	incrCmd := pipe.Incr(ctx, key)

	// Set expiration only if key is new
	pipe.ExpireNX(ctx, key, config.Window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}

	count := int(incrCmd.Val())
	remaining := config.Requests - count
	if remaining < 0 {
		remaining = 0
	}

	resetTime := windowEnd.Unix()

	return count <= config.Requests, remaining, resetTime, nil
}

// EndpointRateLimitMiddleware creates a rate limiter specific to endpoints
func EndpointRateLimitMiddleware(limiter *RateLimiter, limits map[string]RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()

		config, exists := limits[endpoint]
		if !exists {
			// Use default limit
			config = RateLimitConfig{
				Requests: 100,
				Window:   time.Minute,
			}
		}

		identifier := c.ClientIP()
		if userID, exists := c.Get(UserIDKey); exists {
			identifier = userID.(string)
		}

		key := limiter.keyPrefix + ":ratelimit:" + endpoint + ":" + identifier

		allowed, remaining, resetTime, err := limiter.Allow(c.Request.Context(), key, config)
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded for this endpoint",
				"retry_after": resetTime - time.Now().Unix(),
			})
			return
		}

		c.Next()
	}
}

