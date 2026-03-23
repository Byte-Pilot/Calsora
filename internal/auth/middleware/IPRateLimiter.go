package middleware

import (
	"Calsora/pkg/ratelimiter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

func IPRateLimitMiddleware(limiter ratelimiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()
		limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
		if err != nil {
			limit = 60
		}
		window, err := strconv.Atoi(os.Getenv("RATE_WINDOW"))
		if err != nil {
			window = 60
		}

		allowed, err := limiter.LimitByIP(userIP, limit, window)
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
