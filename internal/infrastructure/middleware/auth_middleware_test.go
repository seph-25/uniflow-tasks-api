package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	var capturedUserID string
	
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		require.True(t, exists, "userID should exist in context")
		capturedUserID = userID.(string)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-ID", "user-123")
	req.Header.Set("X-User-Email", "user@uniflow.edu")
	req.Header.Set("X-User-Name", "Juan Perez")
	req.Header.Set("X-User-Picture", "https://example.com/pic.jpg")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "user-123", capturedUserID)
}

func TestAuthMiddleware_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No X-User-ID header
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "UNAUTHORIZED")
}

func TestAuthMiddleware_EmptyUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-ID", "") // Empty string
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_OptionalHeadersMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	var capturedUser interface{}
	
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		user, exists := c.Get("user")
		require.True(t, exists)
		capturedUser = user
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-User-ID", "user-456")
	// Omitir X-User-Email, X-User-Name, X-User-Picture

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotNil(t, capturedUser)
}
