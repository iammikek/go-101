package app

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the API key from the X-API-Key header.
func AuthMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		xAPIKey := c.GetHeader("X-API-Key")
		if xAPIKey == "" || xAPIKey != key {
			c.JSON(401, gin.H{"detail": "Invalid or missing API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
