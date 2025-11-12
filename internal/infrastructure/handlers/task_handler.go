package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"uniflow-api/internal/application"
	"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"
	"uniflow-api/internal/infrastructure/handlers/requests"

	"github.com/gin-gonic/gin"
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

// getUserID extrae el userID del contexto Gin (inyectado por AuthMiddleware)
// Retorna el userID y un bool indicando si se encontró correctamente
func getUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return "", false
	}
	
	return userIDStr, true
}

// GetTasks maneja GET /tasks con filtros opcionales
func (th *TaskHandler) GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}

	// Parsear query parameters
	var filterReq requests.TaskFilterRequest
	if err := c.ShouldBindQuery(&filterReq); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_FILTER", err.Error()))
		return
	}

	// Convertir a domain filter
	filter, err := filterReq.ToTaskFilter(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("INVALID_FILTER", err.Error()))
		return
	}
	if filter == nil {
		filter = &ports.TaskFilter{UserID: userID}
	}

	// Obtener tareas filtradas
	tasks, pageInfo, err := th.taskService.GetTasksFiltered(ctx, *filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = TaskFromDomain(&t)
	}

	response := gin.H{
		"data": taskDTOs,
		"pagination": gin.H{
			"page":       pageInfo.Page,
			"limit":      pageInfo.Limit,
			"total":      pageInfo.Total,
			"totalPages": pageInfo.TotalPages,
			"hasNext":    pageInfo.HasNext,
			"hasPrev":    pageInfo.HasPrev,
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

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
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

// SearchTasks maneja GET /tasks/search
func (th *TaskHandler) SearchTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, NewErrorResponse("MISSING_QUERY", "search query parameter 'q' is required"))
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Crear filter con búsqueda
	filter := ports.TaskFilter{
		UserID:    userID,
		Search:    query,
		Limit:     limit,
		Page:      1,
		SortBy:    "createdAt",
		SortOrder: "desc",
	}

	tasks, pageInfo, err := th.taskService.GetTasksFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("SEARCH_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = TaskFromDomain(&t)
	}

	c.JSON(http.StatusOK, gin.H{
		"query":      query,
		"results":    taskDTOs,
		"count":      len(taskDTOs),
		"totalFound": pageInfo.Total,
	})
}

// GetOverdue maneja GET /tasks/overdue
func (th *TaskHandler) GetOverdue(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}

	tz := c.DefaultQuery("tz", "UTC")

	filter := ports.TaskFilter{
		UserID:    userID,
		IsOverdue: &[]bool{true}[0],
		TimeZone:  tz,
		Limit:     100,
		Page:      1,
		SortBy:    "dueDate",
		SortOrder: "asc",
	}

	tasks, _, err := th.taskService.GetTasksFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		dto := TaskFromDomain(&t)
		// Calcular daysOverdue

		taskDTOs[i] = dto
	}

	c.JSON(http.StatusOK, gin.H{

		"tasks":    taskDTOs,
		"count":    len(taskDTOs),
		"timezone": tz,
	})
}

// GetCompleted maneja GET /tasks/completed
func (th *TaskHandler) GetCompleted(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	filter := ports.TaskFilter{
		UserID:    userID,
		Status:    []string{"done"},
		Limit:     limit,
		Page:      1,
		SortBy:    "updatedAt",
		SortOrder: "desc",
	}

	tasks, pageInfo, err := th.taskService.GetTasksFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = TaskFromDomain(&t)
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":      taskDTOs,
		"count":      len(taskDTOs),
		"pagination": pageInfo,
	})
}

// GetBySubject maneja GET /tasks/by-subject/:subjectId
func (th *TaskHandler) GetBySubject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}
	subjectID := c.Param("subjectId")

	filter := ports.TaskFilter{
		UserID:    userID,
		SubjectID: subjectID,
		Limit:     100,
		Page:      1,
		SortBy:    "dueDate",
		SortOrder: "asc",
	}

	tasks, _, err := th.taskService.GetTasksFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = TaskFromDomain(&t)
	}

	c.JSON(http.StatusOK, gin.H{
		"subjectId": subjectID,
		"tasks":     taskDTOs,
		"count":     len(taskDTOs),
	})
}

// GetByPeriod maneja GET /tasks/by-period/:periodId
func (th *TaskHandler) GetByPeriod(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, NewErrorResponse("UNAUTHORIZED", "userID not found in context (middleware failed)"))
		return
	}
	periodID := c.Param("periodId")

	filter := ports.TaskFilter{
		UserID:    userID,
		PeriodID:  periodID,
		Limit:     100,
		Page:      1,
		SortBy:    "dueDate",
		SortOrder: "asc",
	}

	tasks, _, err := th.taskService.GetTasksFiltered(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("INTERNAL_ERROR", err.Error()))
		return
	}

	taskDTOs := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		taskDTOs[i] = TaskFromDomain(&t)
	}

	c.JSON(http.StatusOK, gin.H{
		"periodId": periodID,
		"tasks":    taskDTOs,
		"count":    len(taskDTOs),
	})
}
