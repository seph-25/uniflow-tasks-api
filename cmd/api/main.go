package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Description string     `json:"description,omitempty"`
}

func main() {
	// Puerto: por defecto 8080. Si existe PORT en env, lo usamos (útil para contenedores).
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// GET /health -> quick check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// GET /tasks -> datos de prueba (mock)
	r.GET("/tasks", func(c *gin.Context) {
		now := time.Now()
		tasks := []Task{
			{ID: "t-001", Title: "Demo: tarea de ejemplo", Status: "todo", Priority: "medium", DueDate: &now, Tags: []string{"demo"}},
			{ID: "t-002", Title: "Otra tarea", Status: "in-progress", Priority: "high"},
		}
		c.JSON(http.StatusOK, tasks)
	})

	// Levanta el server
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
