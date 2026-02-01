package middleware

import (
	jwt2 "Calsora/internal/auth/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		if cookie, err := c.Cookie("access_token"); err == nil {
			tokenStr = cookie
		}

		if tokenStr == "" {
			header := c.GetHeader("access_token")
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		}
		if tokenStr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &jwt2.UserJWTClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET_ACCESS")), nil
		}, jwt.WithLeeway(5*time.Second))
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		claims, ok := token.Claims.(*jwt2.UserJWTClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("subscription", claims.Subscription)
		c.Next()
	}
}
