package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
)

//helper to ensure context is not nil
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// TaskService coordina casos de uso relacionados con tareas
// Ahora depende de una abstracción (TaskRepository) en lugar de datos hardcodeados
type TaskService struct {
	repo        ports.TaskRepository
	queueClient *azqueue.QueueClient
}

// NewTaskService crea una nueva instancia de TaskService
// Inyecta el repositorio (puede ser MongoDB, PostgreSQL, etc.) y opcionalmente un queueClient
func NewTaskService(repo ports.TaskRepository, queueClient *azqueue.QueueClient) *TaskService {
	return &TaskService{
		repo:        repo,
		queueClient: queueClient,
	}
}

// GetAllTasks obtiene todas las tareas del usuario desde la BD real
func (ts *TaskService) GetAllTasks(ctx context.Context, userID string) ([]domain.Task, error) {
	ctx = ensureContext(ctx)
	// Validar que el contexto no fue cancelado
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Llamar al repositorio (ahora MongoDB en lugar de mocks)
	tasks, err := ts.repo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID obtiene una tarea específica por ID
func (ts *TaskService) GetTaskByID(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	task, err := ts.repo.GetByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetTasksByStatus obtiene tareas filtrando por estado
func (ts *TaskService) GetTasksByStatus(ctx context.Context, userID, status string) ([]domain.Task, error) {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tasks, err := ts.repo.GetByUserAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// CreateTask crea una nueva tarea
func (ts *TaskService) CreateTask(ctx context.Context, task *domain.Task, userID, userName, userEmail string) error {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validar que la tarea es válida
	if err := task.IsValid(); err != nil {
		return err
	}

	// Persistir en BD
	err := ts.repo.Create(ctx, task)
	if err != nil {
		return err
	}

	// Encolar mensaje de recordatorio si está disponible Azure Queue
	if ts.queueClient != nil {
		if err := ts.enqueueDeadlineReminder(ctx, task, userID, userName, userEmail); err != nil {
			// Log pero no fallar la creación de tarea
			log.Printf("⚠️ Error al encolar recordatorio para tarea %s: %v", task.ID, err)
		}
	}

	return nil
}

// enqueueDeadlineReminder encola un mensaje para recordatorio de deadline
func (ts *TaskService) enqueueDeadlineReminder(ctx context.Context, task *domain.Task, userID, userName, userEmail string) error {
	// Calcular visibility timeout: 3 días antes del vencimiento
	now := time.Now()
	timeUntilDue := task.DueDate.Sub(now)
	threeDays := 3 * 24 * time.Hour

	var visibilityTimeoutSeconds int32
	if timeUntilDue > threeDays {
		visibilityTimeoutSeconds = int32((timeUntilDue - threeDays).Seconds())
	} else {
		// Si vence en menos de 3 días, hacer visible inmediatamente
		visibilityTimeoutSeconds = 0
	}

	// Construir mensaje JSON
	message := map[string]interface{}{
		"taskId":    task.ID,
		"userId":    userID,
		"name":      userName,
		"email":     userEmail,
		"title":     task.Title,
		"message":   fmt.Sprintf("La tarea '%s' está próxima a vencerse. Faltan 3 días", task.Title),
		"type":      "deadline_reminder",
		"priority":  task.Priority,
		"dueDate":   task.DueDate.Format(time.RFC3339),
		"createdAt": time.Now().Format(time.RFC3339),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error al serializar mensaje: %w", err)
	}

	// Encolar mensaje
	_, err = ts.queueClient.EnqueueMessage(ctx, string(messageJSON), &azqueue.EnqueueMessageOptions{
		VisibilityTimeout: &visibilityTimeoutSeconds,
	})
	if err != nil {
		return fmt.Errorf("error al encolar mensaje: %w", err)
	}

	log.Printf("✅ Recordatorio encolado para tarea %s (visible en %.0f horas)", task.ID, float64(visibilityTimeoutSeconds)/3600)
	return nil
}

// UpdateTask actualiza una tarea existente
func (ts *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validar que puede ser modificada
	if err := task.CanBeModified(); err != nil {
		return err
	}

	// Validar estructura
	if err := task.IsValid(); err != nil {
		return err
	}

	// Persistir cambios
	err := ts.repo.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask elimina una tarea
func (ts *TaskService) DeleteTask(ctx context.Context, taskID, userID string) error {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := ts.repo.Delete(ctx, taskID, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetTasksFiltered obtiene tareas con filtros avanzados
func (ts *TaskService) GetTasksFiltered(ctx context.Context, filter ports.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return nil, domain.PageInfo{}, ctx.Err()
	default:
	}

	tasks, pageInfo, err := ts.repo.FindByFilter(ctx, filter)
	if err != nil {
		return nil, domain.PageInfo{}, err
	}

	return tasks, pageInfo, nil
}

// GetDashboard retorna datos agregados para el dashboard
func (ts *TaskService) GetDashboard(ctx context.Context, userID string) (*domain.DashboardData, error) {
	ctx = ensureContext(ctx)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, err := ts.repo.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
