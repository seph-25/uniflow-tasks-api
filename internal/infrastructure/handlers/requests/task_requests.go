package requests

import (
	"time"
)

// CreateTaskRequest estructura para POST /tasks
type CreateTaskRequest struct {
	Title              string    `json:"title" binding:"required"`
	Description        string    `json:"description"`
	SubjectID          string    `json:"subjectId" binding:"required"`
	PeriodID           string    `json:"periodId"`
	DueDate            time.Time `json:"dueDate" binding:"required"`
	Priority           string    `json:"priority" binding:"required,oneof=low medium high urgent"`
	Type               string    `json:"type" binding:"required,oneof=assignment exam reading presentation lab quiz essay group-work"`
	EstimatedTimeHours int       `json:"estimatedTimeHours"`
	Tags               []string  `json:"tags"`
	IsGroupWork        bool      `json:"isGroupWork"`
	GroupMembers       []string  `json:"groupMembers"`
}

// UpdateTaskRequest estructura para PUT /tasks/:id
type UpdateTaskRequest struct {
	Title              string    `json:"title" binding:"required"`
	Description        string    `json:"description"`
	SubjectID          string    `json:"subjectId" binding:"required"`
	PeriodID           string    `json:"periodId"`
	DueDate            time.Time `json:"dueDate" binding:"required"`
	Priority           string    `json:"priority" binding:"required,oneof=low medium high urgent"`
	Type               string    `json:"type" binding:"required,oneof=assignment exam reading presentation lab quiz essay group-work"`
	EstimatedTimeHours int       `json:"estimatedTimeHours"`
	Tags               []string  `json:"tags"`
	IsGroupWork        bool      `json:"isGroupWork"`
	GroupMembers       []string  `json:"groupMembers"`
}

// UpdateTaskStatusRequest estructura para PATCH /tasks/:id/status
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=todo in-progress in-review done cancelled"`
}

// UpdateTaskCompleteRequest estructura para PATCH /tasks/:id/complete
type UpdateTaskCompleteRequest struct {
	ActualTimeHours int `json:"actualTimeHours"`
}
