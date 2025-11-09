package domain

import (
	"testing"
	"time"
)

func TestTaskIsValid(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &Task{
				Title:     "Test Task",
				SubjectID: "subject-1",
				Status:    StatusTodo,
				Priority:  PriorityMedium,
				Type:      TypeAssignment,
				DueDate:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing title",
			task: &Task{
				Title:     "",
				SubjectID: "subject-1",
				Status:    StatusTodo,
			},
			wantErr: true,
		},
		{
			name: "missing subjectId",
			task: &Task{
				Title:     "Test",
				SubjectID: "",
				Status:    StatusTodo,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.IsValid()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskCanBeModified(t *testing.T) {
	completedTask := &Task{Status: StatusDone}
	if err := completedTask.CanBeModified(); err == nil {
		t.Error("Expected error when modifying completed task")
	}

	openTask := &Task{Status: StatusTodo}
	if err := openTask.CanBeModified(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTaskCanBeDeleted(t *testing.T) {
	completedTask := &Task{Status: StatusDone}
	if err := completedTask.CanBeDeleted(); err == nil {
		t.Error("Expected error when deleting completed task")
	}

	openTask := &Task{Status: StatusTodo}
	if err := openTask.CanBeDeleted(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestTaskIsCompleted(t *testing.T) {
	task := &Task{Status: StatusDone}
	if !task.IsCompleted() {
		t.Error("Expected task to be completed")
	}

	task.Status = StatusTodo
	if task.IsCompleted() {
		t.Error("Expected task not to be completed")
	}
}
