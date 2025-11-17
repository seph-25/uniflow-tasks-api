package memory

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu   sync.RWMutex
	data map[string]*domain.Task
}

func NewRepo() *Repo {
	return &Repo{data: make(map[string]*domain.Task)}
}

func (r *Repo) nextID() string {
	return "t-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (r *Repo) Create(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task.ID == "" {
		task.ID = r.nextID()
	}
	cp := *task
	r.data[task.ID] = &cp
	return nil
}

func (r *Repo) GetByID(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.data[taskID]
	if !ok {
		return nil, ErrNotFound
	}
	// si usás pertenencia por usuario
	if userID != "" && t.UserID != userID {
		return nil, ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *Repo) GetAll(ctx context.Context, userID string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0)
	for _, t := range r.data {
		if userID == "" || t.UserID == userID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *Repo) GetByUserAndStatus(ctx context.Context, userID, status string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Task, 0)
	for _, t := range r.data {
		if (userID == "" || t.UserID == userID) && t.Status == status {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *Repo) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	old, ok := r.data[task.ID]
	if !ok {
		return ErrNotFound
	}
	// pertenencia
	if task.UserID != "" && old.UserID != task.UserID {
		return ErrNotFound
	}
	cp := *task
	r.data[task.ID] = &cp
	return nil
}

func (r *Repo) Delete(ctx context.Context, taskID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.data[taskID]
	if !ok {
		return ErrNotFound
	}
	// pertenencia
	if userID != "" && t.UserID != userID {
		return ErrNotFound
	}
	delete(r.data, taskID)
	return nil
}

// Métodos para Fase 3 (stubs por ahora)
func (r *Repo) Find(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	return []domain.Task{}, domain.PageInfo{}, nil
}

func (r *Repo) DueToday(ctx context.Context, userID string, loc *time.Location) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (r *Repo) Search(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	return []domain.Task{}, domain.PageInfo{}, nil
}

func (r *Repo) Aggregated(ctx context.Context, userID string, until time.Time) (domain.Stats, error) {
	return domain.Stats{}, nil
}

func (r *Repo) FindByFilter(ctx context.Context, filter ports.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Obtener todas las tareas del usuario
	all := make([]domain.Task, 0)
	for _, t := range r.data {
		if t.UserID == filter.UserID {
			all = append(all, *t)
		}
	}

	// APLICAR FILTROS
	filtered := make([]domain.Task, 0)
	for _, t := range all {
		// Filtro por status
		if len(filter.Status) > 0 {
			found := false
			for _, s := range filter.Status {
				if t.Status == s {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por priority
		if len(filter.Priority) > 0 {
			found := false
			for _, p := range filter.Priority {
				if t.Priority == p {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por type
		if len(filter.Type) > 0 {
			found := false
			for _, tp := range filter.Type {
				if t.Type == tp {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Filtro por subject
		if filter.SubjectID != "" && t.SubjectID != filter.SubjectID {
			continue
		}

		// Filtro por period
		if filter.PeriodID != "" && t.PeriodID != filter.PeriodID {
			continue
		}

		// Filtro por fecha vencida (isOverdue)
		if filter.IsOverdue != nil && *filter.IsOverdue {
			if t.DueDate.After(time.Now()) || t.Status == domain.StatusDone {
				continue
			}
		}

		filtered = append(filtered, t)
	}

	// PAGINACIÓN
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	start := (page - 1) * limit
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	slice := filtered[start:end]

	total := int64(len(filtered))
	totalPages := (total + int64(limit) - 1) / int64(limit)
	pi := domain.PageInfo{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    int64(page) < totalPages,
		HasPrev:    page > 1,
	}

	return slice, pi, nil
}
