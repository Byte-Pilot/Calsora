package middleware

import (
	"Calsora/pkg/ratelimiter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func IPRateLimitMiddleware(limiter ratelimiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()

		allowed, err := limiter.LimitByIP(userIP)
		if err != nil {
			log.Printf("rate limiter error: %v", err)
			c.Next()
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Слишком много запросов. Попробуйте позже."})
			return
		}
		c.Next()
		return

	}
}
