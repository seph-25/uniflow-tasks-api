package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/goccy/go-json"
)

// helper to ensure context is not nil
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
// Inyecta el repositorio (puede ser MongoDB, PostgreSQL, etc.)
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

// CreateTask crea una nueva tarea y la encola para recordatorio
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

	// ✅ Encolar mensaje de recordatorio (solo si hay queue client)
	if ts.queueClient != nil {
		if err := ts.enqueueDeadlineReminder(ctx, task, userID, userName, userEmail); err != nil {
			// Log pero no fallar la creación de tarea
			log.Printf("⚠️ Error al encolar recordatorio para tarea %s: %v", task.ID, err)
		}
	}

	return nil
}

// enqueueDeadlineReminder encola un recordatorio 3 días antes del vencimiento
func (ts *TaskService) enqueueDeadlineReminder(ctx context.Context, task *domain.Task, userID, userName, userEmail string) error {
	// 1) Calcular cuándo enviar el recordatorio (3 días antes del vencimiento)
	reminderDays := 3
	now := time.Now()
	reminderTime := task.DueDate.AddDate(0, 0, -reminderDays)

	// 2) Calcular visibilityTimeout: diferencia entre ahora y cuándo se debe procesar
	duration := reminderTime.Sub(now)

	// Si ya pasó la fecha de recordatorio o es muy próxima, enviar inmediatamente
	if duration <= 0 {
		duration = 1 * time.Second
	}

	visibilitySeconds := int32(duration.Seconds())

	// 3) Validar límites de Azure Queue (máximo 7 días = 604800 segundos)
	maxVisibilitySeconds := int32(7 * 24 * 3600)
	if visibilitySeconds > maxVisibilitySeconds {
		visibilitySeconds = maxVisibilitySeconds
		log.Printf("⚠️ Vencimiento de tarea %s excede 7 días - limitado a máximo de Azure", task.ID)
	}

	// 4) Construir payload del mensaje
	data := map[string]interface{}{
		"taskId":    task.ID,
		"userId":    userID,
		"name":      userName,
		"email":     userEmail,
		"title":     task.Title,
		"message":   fmt.Sprintf("La tarea '%s' está próxima a vencerse. Faltan %d días", task.Title, reminderDays),
		"type":      "deadline_reminder",
		"priority":  "high",
		"dueDate":   task.DueDate.Format(time.RFC3339),
		"createdAt": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error serializando mensaje: %w", err)
	}

	// 5) Encolar mensaje
	_, err = ts.queueClient.EnqueueMessage(
		ctx,
		string(jsonData),
		&azqueue.EnqueueMessageOptions{
			VisibilityTimeout: &visibilitySeconds,
		},
	)

	if err != nil {
		return fmt.Errorf("error encolando mensaje: %w", err)
	}

	log.Printf("Recordatorio encolado para tarea %s - Se enviará en %.1f horas (%d segundos)",
		task.ID, duration.Hours(), visibilitySeconds)

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
