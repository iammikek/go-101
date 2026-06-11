package unit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iammikek/go-101/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareValid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/protected", app.AuthMiddleware("test-key-123"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/protected", nil)
	req.Header.Set("X-API-Key", "test-key-123")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddlewareMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/protected", app.AuthMiddleware("test-key-123"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/protected", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddlewareInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/protected", app.AuthMiddleware("test-key-123"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/protected", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
