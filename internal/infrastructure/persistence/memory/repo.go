package memory

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	d "uniflow-api/internal/domain"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu   sync.RWMutex
	data map[string]*d.Task
}

func NewRepo() *Repo {
	return &Repo{data: make(map[string]*d.Task)}
}

func (r *Repo) nextID() string {
	return "t-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func (r *Repo) Create(ctx context.Context, task *d.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task.ID == "" {
		task.ID = r.nextID()
	}
	cp := *task
	r.data[task.ID] = &cp
	return nil
}

func (r *Repo) GetByID(ctx context.Context, taskID, userID string) (*d.Task, error) {
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

func (r *Repo) GetAll(ctx context.Context, userID string) ([]d.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]d.Task, 0)
	for _, t := range r.data {
		if userID == "" || t.UserID == userID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *Repo) GetByUserAndStatus(ctx context.Context, userID, status string) ([]d.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]d.Task, 0)
	for _, t := range r.data {
		if (userID == "" || t.UserID == userID) && t.Status == status {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *Repo) Update(ctx context.Context, task *d.Task) error {
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
func (r *Repo) Find(ctx context.Context, f d.TaskFilter) ([]d.Task, d.PageInfo, error) {
	return []d.Task{}, d.PageInfo{}, nil
}

func (r *Repo) DueToday(ctx context.Context, userID string, loc *time.Location) ([]d.Task, error) {
	return []d.Task{}, nil
}

func (r *Repo) Search(ctx context.Context, f d.TaskFilter) ([]d.Task, d.PageInfo, error) {
	return []d.Task{}, d.PageInfo{}, nil
}

func (r *Repo) Aggregated(ctx context.Context, userID string, until time.Time) (d.Stats, error) {
	return d.Stats{}, nil
}
