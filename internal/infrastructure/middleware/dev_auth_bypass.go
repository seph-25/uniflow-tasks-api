package middleware

import (
	"log"
	"os"

	"uniflow-api/internal/domain"

	"github.com/gin-gonic/gin"
)

// DevAuthBypass permite autenticación local en modo debug
// usando headers X-Dev-User-* en lugar de los inyectados por API Management
//
// SOLO para desarrollo local. NUNCA usar en producción.
//
// Uso:
//   curl -H "X-Dev-User-ID: local-dev-user" \
//        -H "X-Dev-User-Email: dev@uniflow.edu" \
//        http://localhost:8080/tasks
func DevAuthBypass() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo activo en modo debug
		if os.Getenv("GIN_MODE") != "debug" {
			c.Next()
			return
		}

		// Si hay header X-Dev-User-ID, usarlo
		devUserID := c.GetHeader("X-Dev-User-ID")
		if devUserID == "" {
			c.Next()
			return
		}

		log.Printf("⚠️  DEV MODE: Using dev bypass auth for user %s", devUserID)

		// Crear contexto de usuario desde headers de desarrollo
		user := domain.NewUserContext(
			devUserID,
			c.GetHeader("X-Dev-User-Email"),
			c.GetHeader("X-Dev-User-Name"),
			c.GetHeader("X-Dev-User-Picture"),
		)

		// Guardar en contexto (igual que AuthMiddleware)
		c.Set("user", user)
		c.Set("userID", user.ID)

		c.Next()
	}
}
