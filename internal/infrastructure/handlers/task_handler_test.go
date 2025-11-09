package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"uniflow-api/internal/application"
	"uniflow-api/internal/domain"
	"uniflow-api/internal/infrastructure/persistence/memory"
)

// Setup helper para tests
func setupTestRouter() (*gin.Engine, *TaskHandler, *application.TaskService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Middleware para simular auth
	r.Use(func(c *gin.Context) {
		c.Set("userID", "user-test")
		c.Next()
	})

	repo := memory.NewRepo()
	service := application.NewTaskService(repo)
	handler := NewTaskHandler(service)

	return r, handler, service
}

func TestGetTasksWithFilters(t *testing.T) {
	r, handler, service := setupTestRouter()

	// Crear tareas de prueba
	task1 := &domain.Task{
		ID:        "task-1",
		UserID:    "user-test",
		Title:     "Proyecto urgente",
		SubjectID: "subject-ic-6821",
		Status:    domain.StatusTodo,
		Priority:  domain.PriorityHigh,
		DueDate:   time.Now().Add(2 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	service.CreateTask(nil, task1)

	// Registrar ruta
	r.GET("/tasks", handler.GetTasks)

	// Test: Filtrar por prioridad
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks?priority=high", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Validar estructura
	if _, ok := response["data"]; !ok {
		t.Error("Response missing 'data' field")
	}
	if _, ok := response["pagination"]; !ok {
		t.Error("Response missing 'pagination' field")
	}
}

func TestSearchTasks(t *testing.T) {
	r, handler, _ := setupTestRouter()

	r.GET("/tasks/search", handler.SearchTasks)

	// Test sin query
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/search", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing query, got %d", w.Code)
	}

	// Test con query
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/tasks/search?q=proyecto", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
