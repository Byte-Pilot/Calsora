package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireActiveSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkSubPlan := c.GetString("subscription")
		checkSubStatus := c.GetString("sub_status")
		if (checkSubPlan != "premium" && checkSubPlan != "trial") || checkSubStatus != "active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "no subscription"})
			c.Abort()
			return
		}
		c.Next()
	}
}
