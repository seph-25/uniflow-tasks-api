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

// TaskFilter estructura para filtrar tareas en consultas
type TaskFilter struct {
	UserID      string    `form:"userId"` // Obligatorio (viene del JWT)
	Status      []string  `form:"status"`      // Ej: "todo,in-progress"
	Priority    []string  `form:"priority"`    // Ej: "high,urgent"
	SubjectID   string    `form:"subjectId"`
	PeriodID    string    `form:"periodId"`
	DueDateFrom time.Time `form:"dueDateFrom"` // ISO 8601
	DueDateTo   time.Time `form:"dueDateTo"`
	IsOverdue   *bool     `form:"isOverdue"`
	IsDueSoon   *bool     `form:"isDueSoon"`   // Próximas 24h
	Search      string    `form:"search"`      // Búsqueda texto
	SortBy      string    `form:"sortBy"`      // dueDate, priority, status, createdAt
	SortOrder   string    `form:"sortOrder"`   // asc, desc
	Page        int       `form:"page"`
	Limit       int       `form:"limit"`
	TimeZone    string    `form:"tz"`          // Ej: America/Costa_Rica
}

// PageInfo metadatos de paginación
type PageInfo struct {
	Total      int64
	Page       int
	Limit      int
	TotalPages int64
	HasNext    bool
	HasPrev    bool
}
