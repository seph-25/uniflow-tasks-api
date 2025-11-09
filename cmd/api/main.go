package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"uniflow-api/internal/infrastructure/auth"

	"uniflow-api/internal/application"
	ports "uniflow-api/internal/application/ports"
	"uniflow-api/internal/infrastructure/handlers"
	"uniflow-api/internal/infrastructure/persistence"            // Mongo repo
	mem "uniflow-api/internal/infrastructure/persistence/memory" // Repo en memoria
)

func main() {
	// 1) Cargar .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// 2) Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3) Modo Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// 4) Crear el repositorio (memoria por defecto si no hay MONGO_URI)
	mongoURI := os.Getenv("MONGO_URI")
	var repo ports.TaskRepository

	if mongoURI == "" {
		log.Println("MONGO_URI no configurada → usando repositorio EN MEMORIA")
		repo = mem.NewRepo()
	} else {
		log.Println("Inicializando repositorio Mongo…")

		mongoDB := os.Getenv("MONGO_DB")
		if mongoDB == "" {
			mongoDB = "uniflowdb" // default sensato
		}

		// Conectar a Mongo con timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
		if err != nil {
			log.Fatalf("ERROR al conectar a MongoDB: %v", err)
		}
		defer client.Disconnect(context.Background())

		// Ping
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("ERROR al hacer ping a MongoDB: %v", err)
		}
		log.Println("✅ Conectado a MongoDB")

		// Selección de colección y repo
		coll := client.Database(mongoDB).Collection("tasks")
		repo = persistence.NewMongoTaskRepository(coll)
	}

	// 5) Servicio + Router + Handlers
	taskService := application.NewTaskService(repo)
	r := gin.Default()

	taskHandler := handlers.NewTaskHandler(taskService)

	// 6) Rutas

	// 6) Rutas
	r.GET("/health", handlers.HealthHandler) // Sin JWT (health check público)

	// Grupo de rutas protegidas con JWT
	protected := r.Group("/")
	protected.Use(auth.AuthMiddleware())
	{
		// Rutas específicas
		protected.GET("/tasks/search", taskHandler.SearchTasks)
		protected.GET("/tasks/overdue", taskHandler.GetOverdue)
		protected.GET("/tasks/completed", taskHandler.GetCompleted)
		protected.GET("/tasks/by-subject/:subjectId", taskHandler.GetBySubject)
		protected.GET("/tasks/by-period/:periodId", taskHandler.GetByPeriod)

		protected.GET("/tasks", taskHandler.GetTasks)
		protected.GET("/tasks/:id", taskHandler.GetTaskByID)
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)
		protected.PATCH("/tasks/:id/complete", taskHandler.CompleteTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	}

	// Endpoint de testing para generar tokens (SOLO DEV)
	if os.Getenv("GIN_MODE") == "debug" {
		r.POST("/auth/test-token", func(c *gin.Context) {
			token, _ := auth.GenerateToken("user-test-001", "student@uniflow.edu")
			c.JSON(200, gin.H{"token": token})
		})
	}

	// 7) Levantar server
	fmt.Printf("Servidor escuchando en puerto %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("ERROR al levantar servidor: %v", err)
	}
}
