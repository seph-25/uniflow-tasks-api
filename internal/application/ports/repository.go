package ports

import (
	"context"
	"time"
	
	"uniflow-api/internal/domain"
)


// Filtros genéricos para listados/búsquedas
type TaskFilter struct {
	UserID    string    // obligatorio para aislar datos por usuario
	Status    string    // todo | in-progress | done | cancelled (opcional)
	SubjectID string    // opcional
	PeriodID  string    // opcional
	From      *time.Time // ventana inicial (opcional)
	To        *time.Time // ventana final (opcional)
	Query     string    // texto libre (3B) opcional
	Limit     int       // paginación
	Offset    int       // paginación
}

// Metadatos de paginación
type PageInfo struct {
	Total      int64
	Page       int
	Limit      int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}

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
	Find(ctx context.Context, f TaskFilter) ([]domain.Task, PageInfo, error)

	// Ventana de "hoy" según zona horaria
	DueToday(ctx context.Context, userID string, loc *time.Location) ([]domain.Task, error)

	// Búsqueda por texto (puede reutilizar Find si no hay índice text)
	Search(ctx context.Context, f TaskFilter) ([]domain.Task, PageInfo, error)

	// Agregaciones para /stats y /dashboard (3B). Si aún no la usarás,
	// podés devolver valores en cero y null.
	Aggregated(ctx context.Context, userID string, until time.Time) (domain.Stats, error)

}