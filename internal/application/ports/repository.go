package ports

import (
	"context"
	"uniflow-api/internal/domain"
)

// TaskRepository define las operaciones de persistencia para tareas
// Esta interfaz permite que TaskService sea agnóstico de la implementación
// (puede ser MongoDB, PostgreSQL, etc.)
type TaskRepository interface {
	// Create inserta una nueva tarea en la BD
	Create(ctx context.Context, task *domain.Task) error

	// GetByID obtiene una tarea específica por ID y verifica pertenencia al usuario
	GetByID(ctx context.Context, taskID, userID string) (*domain.Task, error)

	// GetAll obtiene todas las tareas de un usuario
	// Implementa contexto y timeout para evitar bloqueos
	GetAll(ctx context.Context, userID string) ([]domain.Task, error)

	// GetByUserAndStatus obtiene tareas filtradas por usuario y estado
	GetByUserAndStatus(ctx context.Context, userID, status string) ([]domain.Task, error)

	// Update actualiza una tarea existente (solo si pertenece al usuario)
	Update(ctx context.Context, task *domain.Task) error

	// Delete elimina una tarea (solo si pertenece al usuario)
	Delete(ctx context.Context, taskID, userID string) error
}