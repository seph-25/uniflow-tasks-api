# UniFlow Tasks Backend API

Sistema de gestión de tareas académicas para estudiantes universitarios.

## 📋 Descripción

**UniFlow** es una plataforma que ayuda a estudiantes a organizar y rastrear sus tareas, 
asignaciones y actividades académicas a través de múltiples materias y períodos escolares.

El **Tasks Backend** es el microservicio central que maneja todas las operaciones 
relacionadas con tareas académicas.

## 🏗️ Arquitectura

El proyecto sigue **Clean Architecture** con tres capas:

- **Domain:** Entidades puras y reglas de negocio
- **Application:** Orquestación de casos de uso
- **Infrastructure:** Handlers HTTP, persistencia, configuración

## 🚀 Fase Actual

**Fase 1A** - Fundación con mocks (Implementando)

- ✅ `GET /health` - Health check
- ✅ `GET /tasks` - Obtener tareas (mock)
- 🔄 Preparación para Fase 1B (integración MongoDB)

## 📦 Stack Tecnológico

- **Lenguaje:** Go 1.22+
- **Framework Web:** Gin v1.11.0
- **Containerización:** Docker multi-stage
- **Infraestructura:** Azure Container Instances
- **Base de Datos:** Azure Cosmos DB for MongoDB (Fase 1B)

## 🔧 Requisitos

- Go 1.22+
- Docker (para containerización)
- Git

## 🏃 Ejecución Local

### Opción 1: Directo
```bash
# Clonar repositorio
git clone <repo-url>
cd uniflow-api

# Descargar dependencias
go mod download

# Ejecutar
go run ./cmd/api/main.go
```

Server escucha en `http://localhost:8080`

### Opción 2: Docker
```bash
# Build
docker build -t uniflow-api:dev .

# Run
docker run -p 8080:8080 -e PORT=8080 uniflow-api:dev
```

## 🧪 Endpoints Fase 1A

### GET /health
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:45Z",
  "version": "1.0.0",
  "service": "tasks"
}
```

### GET /tasks
```bash
curl http://localhost:8080/tasks
```

**Response:**
```json
{
  "data": [
    {
      "id": "task-001",
      "title": "Proyecto Programado I - UniFlow",
      "subjectId": "subject-ic-6821",
      "type": "assignment",
      "status": "in-progress",
      "priority": "high",
      "dueDate": "2025-02-04T00:00:00Z",
      "estimatedTimeHours": 40,
      "tags": ["proyecto", "frontend", "azure"],
      "createdAt": "2025-01-10T00:00:00Z",
      "updatedAt": "2025-01-15T00:00:00Z",
      "completedAt": null
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "totalPages": 1,
    "hasNext": false,
    "hasPrevious": false
  }
}
```

## 📝 Variables de Entorno

Ver `.env.example` para configuración completa.

## 📚 Roadmap

- **Fase 1A** (Actual): Fundación con mocks
- **Fase 1B**: Integración MongoDB
- **Fase 2A**: Pipeline CI/CD GitHub Actions
- **Fase 2B**: CRUD completo
- **Fase 3A**: Consultas avanzadas
- **Fase 3B**: Analytics y Dashboard

## 🔒 Seguridad

- Los endpoints requieren autenticación vía JWT (Fase 2B+)
- Las operaciones se filtran por `userID` extraído del token
- Las imágenes Docker usan `distroless` para mínima superficie de ataque

## 📖 Documentación

- OpenAPI Spec: `UniFlow Tasks Service API.openapi+json.json`
- Plan de Refactorización: `REFACTORING.md` (en desarrollo)

## 👥 Contribuidores

- Dev A: Infraestructura inicial y Dockerfile
- Dev B: Integración MongoDB (Fase 1B)

## 📄 Licencia

MIT License - Ver LICENSE para detalles