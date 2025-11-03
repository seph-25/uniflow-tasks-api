package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"uniflow-api/internal/application"
	"uniflow-api/internal/domain"
	"uniflow-api/internal/infrastructure/handlers/requests"
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	tasks, err := th.taskService.GetAllTasks(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

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

// CreateTask maneja POST /tasks
func (th *TaskHandler) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req requests.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	task := &domain.Task{
		UserID:             userID,
		Title:              req.Title,
		Description:        req.Description,
		SubjectID:          req.SubjectID,
		PeriodID:           req.PeriodID,
		DueDate:            req.DueDate,
		Priority:           req.Priority,
		Type:               req.Type,
		Status:             domain.StatusTodo,
		EstimatedTimeHours: req.EstimatedTimeHours,
		Tags:               req.Tags,
		IsGroupWork:        req.IsGroupWork,
		GroupMembers:       req.GroupMembers,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := th.taskService.CreateTask(ctx, task); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_TASK", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, TaskFromDomain(task))
}

// GetTaskByID maneja GET /tasks/:id
func (th *TaskHandler) GetTaskByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	taskID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	task, err := th.taskService.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, NewErrorResponse("NOT_FOUND", "Tarea no encontrada"))
		return
	}

	c.JSON(http.StatusOK, TaskFromDomain(task))
}

// UpdateTask maneja PUT /tasks/:id
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	taskID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	var req requests.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	// Obtener tarea existente
	task, err := th.taskService.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, NewErrorResponse("NOT_FOUND", "Tarea no encontrada"))
		return
	}

	// Actualizar campos
	task.Title = req.Title
	task.Description = req.Description
	task.SubjectID = req.SubjectID
	task.PeriodID = req.PeriodID
	task.DueDate = req.DueDate
	task.Priority = req.Priority
	task.Type = req.Type
	task.EstimatedTimeHours = req.EstimatedTimeHours
	task.Tags = req.Tags
	task.IsGroupWork = req.IsGroupWork
	task.GroupMembers = req.GroupMembers
	task.UpdatedAt = time.Now()

	if err := th.taskService.UpdateTask(ctx, task); err != nil {
		c.JSON(http.StatusConflict, NewErrorResponse("CONFLICT", err.Error()))
		return
	}

	c.JSON(http.StatusOK, TaskFromDomain(task))
}

// UpdateTaskStatus maneja PATCH /tasks/:id/status
func (th *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	taskID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	var req requests.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_REQUEST", err.Error()))
		return
	}

	task, err := th.taskService.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, NewErrorResponse("NOT_FOUND", "Tarea no encontrada"))
		return
	}

	task.Status = req.Status
	task.UpdatedAt = time.Now()

	if req.Status == domain.StatusDone && task.CompletedAt == nil {
		now := time.Now()
		task.CompletedAt = &now
	}

	if err := th.taskService.UpdateTask(ctx, task); err != nil {
		c.JSON(http.StatusConflict, NewErrorResponse("CONFLICT", err.Error()))
		return
	}

	c.JSON(http.StatusOK, TaskFromDomain(task))
}

// DeleteTask maneja DELETE /tasks/:id
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	taskID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	if err := th.taskService.DeleteTask(ctx, taskID, userID); err != nil {
		c.JSON(http.StatusConflict, NewErrorResponse("CONFLICT", err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// CompleteTask maneja PATCH /tasks/:id/complete
func (th *TaskHandler) CompleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	taskID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "user-demo"
	}

	var req requests.UpdateTaskCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.ActualTimeHours = 0
	}

	task, err := th.taskService.GetTaskByID(ctx, taskID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, NewErrorResponse("NOT_FOUND", "Tarea no encontrada"))
		return
	}

	task.Status = domain.StatusDone
	if req.ActualTimeHours > 0 {
		task.ActualTimeHours = &[]int{req.ActualTimeHours}[0]
	}
	now := time.Now()
	task.CompletedAt = &now
	task.UpdatedAt = now

	if err := th.taskService.UpdateTask(ctx, task); err != nil {
		c.JSON(http.StatusConflict, NewErrorResponse("CONFLICT", err.Error()))
		return
	}

	c.JSON(http.StatusOK, TaskFromDomain(task))
}