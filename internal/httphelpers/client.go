package httphelpers

import "github.com/gin-gonic/gin"

func IsWebClient(c *gin.Context) bool {
	return c.GetHeader("X-Client-Type") == "web"
}
