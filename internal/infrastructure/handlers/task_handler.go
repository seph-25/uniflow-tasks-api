package handlers

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "uniflow-api/internal/application"
)

// TaskHandler maneja operaciones de tareas
type TaskHandler struct {
    taskService *application.TaskService
}

// NewTaskHandler crea un nuevo TaskHandler
func NewTaskHandler(ts *application.TaskService) *TaskHandler {
    return &TaskHandler{
        taskService: ts,
    }
}

// GetTasks maneja GET /tasks
func (th *TaskHandler) GetTasks(c *gin.Context) {
    // Crear contexto con timeout de 5 segundos
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()

    // En Fase 1B, el userID vendrá del JWT del middleware de autenticación
    // Por ahora usamos un valor de prueba
    userID := c.GetString("userID")
    if userID == "" {
        userID = "user-demo"
    }

    tasks, err := th.taskService.GetAllTasks(ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
        return
    }

    // Convertir tasks a DTOs
    taskDTOs := make([]TaskDTO, len(tasks))
    for i, t := range tasks {
        taskDTOs[i] = TaskFromDomain(&t)
    }

    response := GetTasksResponse{
        Data: taskDTOs,
        Pagination: Pagination{
            Page:        1,
            Limit:       10,
            Total:       len(taskDTOs),
            TotalPages:  1,
            HasNext:     false,
            HasPrevious: false,
        },
    }

    c.JSON(http.StatusOK, response)
}