package application

import (
	"context"
	"testing"
	//"uniflow-api/internal/application/ports"
	"uniflow-api/internal/domain"
)

// Mock repository para tests
type mockRepository struct {
	tasks map[string][]domain.Task
}

func (m *mockRepository) Create(ctx context.Context, task *domain.Task) error {
	if m.tasks == nil {
		m.tasks = make(map[string][]domain.Task)
	}
	m.tasks[task.UserID] = append(m.tasks[task.UserID], *task)
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	return nil, nil // Implementar si necesario
}

func (m *mockRepository) GetAll(ctx context.Context, userID string) ([]domain.Task, error) {
	if tasks, ok := m.tasks[userID]; ok {
		return tasks, nil
	}
	return []domain.Task{}, nil
}

func (m *mockRepository) GetByUserAndStatus(ctx context.Context, userID, status string) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (m *mockRepository) Update(ctx context.Context, task *domain.Task) error {
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, taskID, userID string) error {
	return nil
}

func TestGetAllTasks(t *testing.T) {
	repo := &mockRepository{
		tasks: map[string][]domain.Task{
			"user-1": {
				{
					ID:        "task-1",
					UserID:    "user-1",
					Title:     "Task 1",
					SubjectID: "subject-1",
					Status:    domain.StatusTodo,
				},
			},
		},
	}

	service := NewTaskService(repo)
	tasks, err := service.GetAllTasks(context.Background(), "user-1")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestCreateTask(t *testing.T) {
	repo := &mockRepository{}
	service := NewTaskService(repo)

	task := &domain.Task{
		Title:     "New Task",
		SubjectID: "subject-1",
		Status:    domain.StatusTodo,
		Priority:  domain.PriorityMedium,
		Type:      domain.TypeAssignment,
		UserID:    "user-1",
	}

	err := service.CreateTask(context.Background(), task)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verificar que se guardó
	tasks, _ := repo.GetAll(context.Background(), "user-1")
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task in repository, got %d", len(tasks))
	}
}