package requests

import "strings"
import "time"

import "uniflow-api/internal/application/ports"

// TaskFilterRequest estructura para parsear query parameters
type TaskFilterRequest struct {
	Status      string `form:"status"`       // Comma-separated: "todo,in-progress"
	Priority    string `form:"priority"`     // Comma-separated: "high,urgent"
	SubjectID   string `form:"subjectId"`
	PeriodID    string `form:"periodId"`
	DueDateFrom string `form:"dueDateFrom"`  // ISO 8601: "2025-10-01"
	DueDateTo   string `form:"dueDateTo"`    // ISO 8601: "2025-10-31"
	IsOverdue   *bool  `form:"isOverdue"`    // true/false
	IsDueSoon   *bool  `form:"isDueSoon"`    // true/false (próximas 24h)
	Search      string `form:"search"`       // Búsqueda libre
	SortBy      string `form:"sortBy"`       // dueDate, priority, status, createdAt
	SortOrder   string `form:"sortOrder"`    // asc, desc
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	TimeZone    string `form:"tz"`           // Ej: "America/Costa_Rica"
}

// ToTaskFilter convierte request a domain.TaskFilter
func (req *TaskFilterRequest) ToTaskFilter(userID string) (*ports.TaskFilter, error) {
	filter := &ports.TaskFilter{
		UserID:    userID,
		SubjectID: req.SubjectID,
		PeriodID:  req.PeriodID,
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		IsOverdue: req.IsOverdue,
		IsDueSoon: req.IsDueSoon,
		Page:      req.Page,
		Limit:     req.Limit,
		TimeZone:  req.TimeZone,
	}

	// Parse dates
	if req.DueDateFrom != "" {
		t, err := time.Parse("2006-01-02", req.DueDateFrom)
		if err != nil {
			return nil, err
		}
		filter.DueDateFrom = t
	}

	if req.DueDateTo != "" {
		t, err := time.Parse("2006-01-02", req.DueDateTo)
		if err != nil {
			return nil, err
		}
		filter.DueDateTo = t
	}

	// Parse status (comma-separated)
	if req.Status != "" {
		filter.Status = strings.Split(req.Status, ",")
		for i := range filter.Status {
			filter.Status[i] = strings.TrimSpace(filter.Status[i])
		}
	}

	// Parse priority (comma-separated)
	if req.Priority != "" {
		filter.Priority = strings.Split(req.Priority, ",")
		for i := range filter.Priority {
			filter.Priority[i] = strings.TrimSpace(filter.Priority[i])
		}
	}

	// Validar timezone
	if filter.TimeZone == "" {
		filter.TimeZone = "UTC"
	}

	// Defaults
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Max 100 por página
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	if filter.SortOrder == "" {
		filter.SortOrder = "asc"
	}

	if filter.SortBy == "" {
		filter.SortBy = "dueDate"
	}

	return filter, nil
}

// FilterRequest alias para request generator
type FilterRequest = TaskFilterRequest