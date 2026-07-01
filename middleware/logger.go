package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func statusColor(code int) string {
	switch {
	case code >= 500:
		return "\033[31m" // red
	case code >= 400:
		return "\033[33m" // yellow
	case code >= 300:
		return "\033[36m" // cyan
	default:
		return "\033[32m" // green
	}
}

const reset = "\033[0m"

func Logger(log *zap.Logger, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()

		if env == "production" {
			log.Info("request",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.Int("status", status),
				zap.Duration("latency", time.Since(start)),
			)
		} else {
			log.Info(fmt.Sprintf("%s%d%s %s %s",
				statusColor(status), status, reset,
				c.Request.Method,
				c.Request.URL.Path,
			),
				zap.String("ip", c.ClientIP()),
				zap.Duration("latency", time.Since(start)),
			)
		}
	}
}

func NewLogger(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
