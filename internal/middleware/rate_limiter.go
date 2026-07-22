package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
	prefix string
}

func NewRedisRateLimiter(client *redis.Client, limit int, window time.Duration, prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
		prefix: prefix,
	}
}

func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := rl.prefix + ":" + ip

		ctx := c.Request.Context()

		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			// Redis down — fail open so an outage doesn't take down the whole API.
			// Log it; don't block the request on this.
			c.Next()
			return
		}

		if count == 1 {
			rl.client.Expire(ctx, key, rl.window)
		}

		if count > int64(rl.limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}
