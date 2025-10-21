package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// HealthHandler maneja GET /health
func HealthHandler(c *gin.Context) {
    response := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z07:00"),
        Version:   "1.0.0",
        Service:   "tasks",
    }
    c.JSON(http.StatusOK, response)
}