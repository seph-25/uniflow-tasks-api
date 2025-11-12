package middleware

import (
	"log"
	"net/http"

	"uniflow-api/internal/domain"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware extrae y valida los headers de autenticación
// inyectados por API Management (Google OAuth)
//
// Headers esperados:
//   - X-User-ID (obligatorio)
//   - X-User-Email (opcional)
//   - X-User-Name (opcional)
//   - X-User-Picture (opcional)
//
// Si X-User-ID no está presente, retorna 401 Unauthorized
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer usuario desde headers
		user, err := domain.FromHeaders(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "missing or invalid authentication headers",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Loggear advertencias si headers opcionales faltan
		if user.Email == "" {
			log.Printf("Warning: X-User-Email header missing for user %s", user.ID)
		}
		if user.Name == "" {
			log.Printf("Warning: X-User-Name header missing for user %s", user.ID)
		}

		// Guardar en contexto de Gin para uso en handlers
		c.Set("user", user)
		c.Set("userID", user.ID)

		// Continuar con el siguiente handler
		c.Next()
	}
}
