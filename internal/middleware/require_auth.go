package middleware

import (
	"Calsora/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &utils.UserJWTClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			log.Println(err)
			return
		}

		claims, ok := token.Claims.(*utils.UserJWTClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("subscription", claims.Subscription)
		c.Next()
	}
}
