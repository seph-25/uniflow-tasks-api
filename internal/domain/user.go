package domain

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// UserContext representa el contexto de usuario autenticado
// extraído desde headers HTTP (inyectados por API Management)
type UserContext struct {
	ID      string // X-User-ID (obligatorio)
	Email   string // X-User-Email
	Name    string // X-User-Name
	Picture string // X-User-Picture
}

// NewUserContext crea una nueva instancia de UserContext
func NewUserContext(id, email, name, picture string) *UserContext {
	return &UserContext{
		ID:      id,
		Email:   email,
		Name:    name,
		Picture: picture,
	}
}

// IsValid valida que el contexto de usuario tenga al menos el ID
func (u *UserContext) IsValid() error {
	if u.ID == "" {
		return errors.New("user ID is required")
	}
	return nil
}

// FromHeaders extrae el UserContext desde los headers HTTP del contexto Gin
// Headers esperados: X-User-ID, X-User-Email, X-User-Name, X-User-Picture
// Solo X-User-ID es obligatorio
func FromHeaders(c *gin.Context) (*UserContext, error) {
	userID := c.GetHeader("X-User-ID")
	
	if userID == "" {
		return nil, errors.New("X-User-ID header is required")
	}

	user := &UserContext{
		ID:      userID,
		Email:   c.GetHeader("X-User-Email"),
		Name:    c.GetHeader("X-User-Name"),
		Picture: c.GetHeader("X-User-Picture"),
	}

	if err := user.IsValid(); err != nil {
		return nil, err
	}

	return user, nil
}
