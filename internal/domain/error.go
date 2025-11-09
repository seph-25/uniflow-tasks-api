package domain

import (
	"fmt"
)

// DomainError representa errores del dominio de negocio
type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Errores comunes
var (
	ErrTaskNotFound         = &DomainError{Code: "TASK_NOT_FOUND", Message: "tarea no encontrada"}
	ErrTaskAlreadyCompleted = &DomainError{Code: "TASK_COMPLETED", Message: "la tarea ya está completada"}
	ErrTaskCancelled        = &DomainError{Code: "TASK_CANCELLED", Message: "la tarea está cancelada"}
	ErrInvalidTaskData      = &DomainError{Code: "INVALID_TASK", Message: "datos de tarea inválidos"}
	ErrUnauthorized         = &DomainError{Code: "UNAUTHORIZED", Message: "no autorizado"}
)
