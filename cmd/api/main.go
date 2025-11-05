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

	"uniflow-api/internal/application"
	ports "uniflow-api/internal/application/ports"
	"uniflow-api/internal/infrastructure/handlers"
	"uniflow-api/internal/infrastructure/persistence"          // Mongo repo
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
	r.GET("/health", handlers.HealthHandler)

	r.GET("/tasks", taskHandler.GetTasks)
	r.GET("/tasks/:id", taskHandler.GetTaskByID)
	r.POST("/tasks", taskHandler.CreateTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)
	r.PATCH("/tasks/:id/complete", taskHandler.CompleteTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// 7) Levantar server
	fmt.Printf("Servidor escuchando en puerto %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("ERROR al levantar servidor: %v", err)
	}
}
