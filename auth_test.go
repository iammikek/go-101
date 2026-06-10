package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareValid(t *testing.T) {
	originalAPIKey := apiKey
	apiKey = "test-key-123"
	defer func() { apiKey = originalAPIKey }()

	router := gin.Default()
	router.DELETE("/protected", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/protected", nil)
	req.Header.Set("X-API-Key", "test-key-123")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddlewareMissing(t *testing.T) {
	originalAPIKey := apiKey
	apiKey = "test-key-123"
	defer func() { apiKey = originalAPIKey }()

	router := gin.Default()
	router.DELETE("/protected", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/protected", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddlewareInvalid(t *testing.T) {
	originalAPIKey := apiKey
	apiKey = "test-key-123"
	defer func() { apiKey = originalAPIKey }()

	router := gin.Default()
	router.DELETE("/protected", authMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/protected", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddlewareInit(t *testing.T) {
	os.Setenv("API_KEY", "custom-key")
	defer os.Unsetenv("API_KEY")

	// Re-initialize apiKey
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "dev-key-123"
	}

	assert.Equal(t, "custom-key", apiKey)
}
