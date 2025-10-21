package main

import (
    "os"

    "github.com/gin-gonic/gin"
    "uniflow-api/internal/application"
    "uniflow-api/internal/infrastructure/handlers"
)

func main() {
    // Configurar puerto
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Configurar modo Gin
    if os.Getenv("GIN_MODE") == "" {
        gin.SetMode(gin.DebugMode)
    }

    // Inicializar servicios
    taskService := application.NewTaskService()

    // Crear router
    r := gin.Default()

    // Crear handlers
    taskHandler := handlers.NewTaskHandler(taskService)

    // Rutas
    r.GET("/health", handlers.HealthHandler)
    r.GET("/tasks", taskHandler.GetTasks)

    // Levantar servidor
    if err := r.Run(":" + port); err != nil {
        panic(err)
    }
}