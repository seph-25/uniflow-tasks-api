package domain

import "time"

// Task Status
const (
    StatusTodo       = "todo"
    StatusInProgress = "in-progress"
    StatusInReview   = "in-review"
    StatusDone       = "done"
    StatusCancelled  = "cancelled"
)

var ValidStatuses = []string{StatusTodo, StatusInProgress, StatusInReview, StatusDone, StatusCancelled}

// Task Priority
const (
    PriorityLow    = "low"
    PriorityMedium = "medium"
    PriorityHigh   = "high"
    PriorityUrgent = "urgent"
)

var ValidPriorities = []string{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}

// Task Type
const (
    TypeAssignment = "assignment"
    TypeExam       = "exam"
    TypeReading    = "reading"
    TypePresentation = "presentation"
    TypeLab        = "lab"
    TypeQuiz       = "quiz"
    TypeEssay      = "essay"
    TypeGroupWork  = "group-work"
)

var ValidTypes = []string{TypeAssignment, TypeExam, TypeReading, TypePresentation, TypeLab, TypeQuiz, TypeEssay, TypeGroupWork}

// Stats representa estadísticas agregadas de tareas (para Fase 3B)
type Stats struct {
	Total       int            `json:"total"`
	Completed   int            `json:"completed"`
	Pending     int            `json:"pending"`
	Overdue     int            `json:"overdue"`
	ByStatus    map[string]int `json:"byStatus"`
	ByPriority  map[string]int `json:"byPriority"`
	ByType      map[string]int `json:"byType"`
	BySubject   map[string]int `json:"bySubject"`
}

// TaskFilter filtros genéricos para listados/búsquedas (Fase 3A)
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

// PageInfo metadatos de paginación (Fase 3A)
type PageInfo struct {
	Total      int64
	Page       int
	Limit      int
	TotalPages int
	HasNext    bool
	HasPrev    bool
}