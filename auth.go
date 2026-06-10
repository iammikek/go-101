package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

var apiKey string

func init() {
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "dev-key-123"
	}
}

// authMiddleware verifies the API key from X-API-Key header
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		xAPIKey := c.GetHeader("X-API-Key")
		if xAPIKey == "" || xAPIKey != apiKey {
			c.JSON(401, gin.H{"detail": "Invalid or missing API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
