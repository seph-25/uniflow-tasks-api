package ports

import (
	"context"
	"time"

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

	// Listado con filtros y paginación
	Find(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error)

	// Ventana de "hoy" según zona horaria
	DueToday(ctx context.Context, userID string, loc *time.Location) ([]domain.Task, error)

	// Búsqueda por texto (puede reutilizar Find si no hay índice text)
	Search(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error)

	// Agregaciones para /stats y /dashboard (3B). Si aún no la usarás,
	// podés devolver valores en cero y null.
	Aggregated(ctx context.Context, userID string, until time.Time) (domain.Stats, error)

	// GetDashboardStats retorna estadísticas para el dashboard
	GetDashboardStats(ctx context.Context, userID string) (domain.DashboardData, error)

	FindByFilter(ctx context.Context, filter TaskFilter) ([]domain.Task, domain.PageInfo, error)
}

// TaskFilter es un alias a domain.TaskFilter para mantener compatibilidad
type TaskFilter = domain.TaskFilter
