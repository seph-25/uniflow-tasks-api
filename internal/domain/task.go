package domain

import (
    "time"
    "fmt"
)

// Task representa una actividad académica individual
type Task struct {
    ID                 string     `bson:"_id,omitempty" json:"id"`
    UserID             string     `bson:"userId" json:"-"`                    // No se expone en API
    Title              string     `bson:"title" json:"title"`
    Description        string     `bson:"description" json:"description"`
    SubjectID          string     `bson:"subjectId" json:"subjectId"`         // Materia
    PeriodID           string     `bson:"periodId" json:"periodId"`           // Semestre
    DueDate            time.Time  `bson:"dueDate" json:"dueDate"`
    Status             string     `bson:"status" json:"status"`
    Priority           string     `bson:"priority" json:"priority"`
    Type               string     `bson:"type" json:"type"`
    EstimatedTimeHours int        `bson:"estimatedTimeHours" json:"estimatedTimeHours"`
    ActualTimeHours    *int       `bson:"actualTimeHours,omitempty" json:"actualTimeHours,omitempty"`
    Tags               []string   `bson:"tags" json:"tags"`
    IsGroupWork        bool       `bson:"isGroupWork" json:"isGroupWork"`
    GroupMembers       []string   `bson:"groupMembers" json:"groupMembers"`
    Attachments        []string   `bson:"attachments" json:"attachments"`
    CreatedAt          time.Time  `bson:"createdAt" json:"createdAt"`
    UpdatedAt          time.Time  `bson:"updatedAt" json:"updatedAt"`
    CompletedAt        *time.Time `bson:"completedAt,omitempty" json:"completedAt,omitempty"`
}

// IsValid valida que la Task cumple con reglas de negocio
func (t *Task) IsValid() error {
    if t.Title == "" {
        return fmt.Errorf("título es requerido")
    }
    if t.SubjectID == "" {
        return fmt.Errorf("subjectId es requerido")
    }
    if !isValidStatus(t.Status) {
        return fmt.Errorf("estado inválido: %s", t.Status)
    }
    if !isValidPriority(t.Priority) {
        return fmt.Errorf("prioridad inválida: %s", t.Priority)
    }
    if !isValidType(t.Type) {
        return fmt.Errorf("tipo inválido: %s", t.Type)
    }
    return nil
}

// IsCompleted devuelve si la tarea está completada
func (t *Task) IsCompleted() bool {
    return t.Status == StatusDone
}

// IsCancelled devuelve si la tarea fue cancelada
func (t *Task) IsCancelled() bool {
    return t.Status == StatusCancelled
}

// CanBeModified valida si la tarea puede ser modificada
func (t *Task) CanBeModified() error {
    if t.IsCompleted() {
        return fmt.Errorf("no se puede modificar una tarea completada")
    }
    if t.IsCancelled() {
        return fmt.Errorf("no se puede modificar una tarea cancelada")
    }
    return nil
}

// CanBeDeleted valida si la tarea puede ser eliminada
func (t *Task) CanBeDeleted() error {
    if t.IsCompleted() {
        return fmt.Errorf("no se puede eliminar una tarea completada")
    }
    return nil
}

// Helpers
func isValidStatus(s string) bool {
    for _, v := range ValidStatuses {
        if v == s {
            return true
        }
    }
    return false
}

func isValidPriority(p string) bool {
    for _, v := range ValidPriorities {
        if v == p {
            return true
        }
    }
    return false
}

func isValidType(t string) bool {
    for _, v := range ValidTypes {
        if v == t {
            return true
        }
    }
    return false
}