package application

import (
    "context"
    "time"
    "uniflow-api/internal/domain"
)

// TaskService coordina casos de uso relacionados con tareas
type TaskService struct {
    // En Fase 1A: sin repositorio (mocks)
    // En Fase 1B: aquí iría el TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService() *TaskService {
    return &TaskService{}
}

// GetAllTasks retorna todas las tareas del usuario (mock para Fase 1A)
func (ts *TaskService) GetAllTasks(ctx context.Context, userID string) ([]domain.Task, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // Mock data para Fase 1A
    now := time.Now()
    tasks := []domain.Task{
        {
            ID:                 "task-001",
            UserID:             userID,
            Title:              "Proyecto Programado I - UniFlow",
            Description:        "Desarrollar aplicación web con React y Azure",
            SubjectID:          "subject-ic-6821",
            PeriodID:           "period-2025-01",
            DueDate:            now.AddDate(0, 0, 20),
            Status:             domain.StatusInProgress,
            Priority:           domain.PriorityHigh,
            Type:               domain.TypeAssignment,
            EstimatedTimeHours: 40,
            Tags:               []string{"proyecto", "frontend", "azure"},
            IsGroupWork:        false,
            GroupMembers:       []string{},
            Attachments:        []string{},
            CreatedAt:          now.AddDate(0, 0, -5),
            UpdatedAt:          now,
            CompletedAt:        nil,
        },
        {
            ID:                 "task-002",
            UserID:             userID,
            Title:              "Laboratorio 3 - Consultas SQL",
            Description:        "Implementar consultas con JOIN y subconsultas",
            SubjectID:          "subject-ic-4302",
            PeriodID:           "period-2025-01",
            DueDate:            now.AddDate(0, 0, 10),
            Status:             domain.StatusTodo,
            Priority:           domain.PriorityMedium,
            Type:               domain.TypeLab,
            EstimatedTimeHours: 3,
            Tags:               []string{"lab", "sql", "database"},
            IsGroupWork:        false,
            GroupMembers:       []string{},
            Attachments:        []string{},
            CreatedAt:          now.AddDate(0, 0, -3),
            UpdatedAt:          now.AddDate(0, 0, -3),
            CompletedAt:        nil,
        },
    }

    return tasks, nil
}