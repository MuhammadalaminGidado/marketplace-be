// internal/redis/redis.go
package redis

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func New(addr, password string, db int) *Client {
	// Check for REDIS_URL (Render's format)
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("Failed to parse REDIS_URL: %v, falling back to individual config", err)
		} else {
			client := redis.NewClient(opt)
			return &Client{Client: client}
		}
	}

	// Fallback to individual config for local development
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{Client: client}
}
