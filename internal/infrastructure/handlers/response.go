package handlers

import (
    "uniflow-api/internal/domain"
)

// TaskDTO es la representación de Task en respuestas HTTP
type TaskDTO struct {
    ID                 string         `json:"id"`
    Title              string         `json:"title"`
    Description        string         `json:"description"`
    SubjectID          string         `json:"subjectId"`
    PeriodID           string         `json:"periodId"`
    DueDate            string         `json:"dueDate"`
    Status             string         `json:"status"`
    Priority           string         `json:"priority"`
    Type               string         `json:"type"`
    EstimatedTimeHours int            `json:"estimatedTimeHours"`
    ActualTimeHours    *int           `json:"actualTimeHours,omitempty"`
    Tags               []string       `json:"tags"`
    IsGroupWork        bool           `json:"isGroupWork"`
    GroupMembers       []string       `json:"groupMembers"`
    Attachments        []string       `json:"attachments"`
    CreatedAt          string         `json:"createdAt"`
    UpdatedAt          string         `json:"updatedAt"`
    CompletedAt        *string        `json:"completedAt,omitempty"`
}

// FromDomain convierte domain.Task a TaskDTO
func TaskFromDomain(t *domain.Task) TaskDTO {
    dto := TaskDTO{
        ID:                 t.ID,
        Title:              t.Title,
        Description:        t.Description,
        SubjectID:          t.SubjectID,
        PeriodID:           t.PeriodID,
        DueDate:            t.DueDate.Format("2006-01-02T15:04:05Z07:00"),
        Status:             t.Status,
        Priority:           t.Priority,
        Type:               t.Type,
        EstimatedTimeHours: t.EstimatedTimeHours,
        ActualTimeHours:    t.ActualTimeHours,
        Tags:               t.Tags,
        IsGroupWork:        t.IsGroupWork,
        GroupMembers:       t.GroupMembers,
        Attachments:        t.Attachments,
        CreatedAt:          t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:          t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
    if t.CompletedAt != nil {
        completedStr := t.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
        dto.CompletedAt = &completedStr
    }
    return dto
}

// GetTasksResponse estructura de respuesta para GET /tasks
type GetTasksResponse struct {
    Data       []TaskDTO  `json:"data"`
    Pagination Pagination `json:"pagination"`
}

// Pagination información de paginación
type Pagination struct {
    Page         int `json:"page"`
    Limit        int `json:"limit"`
    Total        int `json:"total"`
    TotalPages   int `json:"totalPages"`
    HasNext      bool `json:"hasNext"`
    HasPrevious  bool `json:"hasPrevious"`
}

// HealthResponse estructura de respuesta para GET /health
type HealthResponse struct {
    Status    string `json:"status"`
    Timestamp string `json:"timestamp"`
    Version   string `json:"version"`
    Service   string `json:"service"`
}

// ErrorResponse estructura de respuesta de error
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// NewErrorResponse crea una respuesta de error
func NewErrorResponse(code, message string) ErrorResponse {
    return ErrorResponse{
        Code:    code,
        Message: message,
    }
}