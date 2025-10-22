package application

import (
	"context"
	"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"
)

// TaskService coordina casos de uso relacionados con tareas
// Ahora depende de una abstracción (TaskRepository) en lugar de datos hardcodeados
type TaskService struct {
	repo ports.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
// Inyecta el repositorio (puede ser MongoDB, PostgreSQL, etc.)
func NewTaskService(repo ports.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// GetAllTasks obtiene todas las tareas del usuario desde la BD real
func (ts *TaskService) GetAllTasks(ctx context.Context, userID string) ([]domain.Task, error) {
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
func (ts *TaskService) CreateTask(ctx context.Context, task *domain.Task) error {
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

	return nil
}

// UpdateTask actualiza una tarea existente
func (ts *TaskService) UpdateTask(ctx context.Context, task *domain.Task) error {
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