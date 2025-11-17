package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"uniflow-api/internal/application"
	ports "uniflow-api/internal/application/ports"
	"uniflow-api/internal/infrastructure/handlers"
	"uniflow-api/internal/infrastructure/middleware"
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

	queueClient := initQueueClient()

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
		defer func() {
			_ = client.Disconnect(context.Background())
		}()

		// Ping
		if err := client.Ping(ctx, nil); err != nil {
			log.Printf("Mongo URI: %v", mongoURI)
			log.Fatalf("ERROR al hacer ping a MongoDB: %v", err)
		}
		log.Println("Conectado a MongoDB")

		// Selección de colección y repo
		coll := client.Database(mongoDB).Collection("tasks")
		repo = persistence.NewMongoTaskRepository(coll)
	}

	// 5) Servicio + Router + Handlers
	taskService := application.NewTaskService(repo, queueClient)
	r := gin.Default()

	taskHandler := handlers.NewTaskHandler(taskService)

	// 6) Rutas públicas (sin autenticación)
	r.GET("/health", handlers.HealthHandler)

	// 7) Middleware de autenticación (headers de API Management)
	// En desarrollo, DevAuthBypass permite usar X-Dev-User-ID
	if os.Getenv("GIN_MODE") == "debug" {
		r.Use(middleware.DevAuthBypass())
	}
	r.Use(middleware.AuthMiddleware())

	// 8) Rutas protegidas (requieren headers X-User-*)
	// Rutas específicas (deben ir primero para no colisionar con :id)
	r.GET("/tasks/search", taskHandler.SearchTasks)
	r.GET("/tasks/overdue", taskHandler.GetOverdue)
	r.GET("/tasks/completed", taskHandler.GetCompleted)
	r.GET("/tasks/by-subject/:subjectId", taskHandler.GetBySubject)
	r.GET("/tasks/by-period/:periodId", taskHandler.GetByPeriod)

	// Rutas CRUD
	r.GET("/tasks", taskHandler.GetTasks)
	r.GET("/tasks/:id", taskHandler.GetTaskByID)
	r.POST("/tasks", taskHandler.CreateTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)
	r.PATCH("/tasks/:id/complete", taskHandler.CompleteTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// 9) Levantar server
	fmt.Printf("Servidor escuchando en puerto %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("ERROR al levantar servidor: %v", err)
	}
}

// QueueClient singleton con lazy initialization
var (
	queueClient *azqueue.QueueClient
	onceQueue   sync.Once
)

func initQueueClient() *azqueue.QueueClient {
	onceQueue.Do(func() {
		queueClient = createQueueConnection()
		log.Println("QueueClient de Azure inicializado")
	})
	return queueClient
}

func createQueueConnection() *azqueue.QueueClient {
	connectionString := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")

	if connectionString == "" {
		log.Println("AZURE_STORAGE_CONNECTION_STRING no configurada - QueueClient no disponible")
		return nil
	}

	queueName := os.Getenv("AZURE_STORAGE_QUEUE_NAME")
	if queueName == "" {
		log.Println("AZURE_STORAGE_QUEUE_NAME no configurada")
		return nil
	}

	queueClient, err := azqueue.NewQueueClientFromConnectionString(
		connectionString,
		queueName,
		nil,
	)
	if err != nil {
		log.Printf("Error al crear queue client: %v", err)
		return nil
	}

	log.Printf("QueueClient creado para cola: %s", queueName)
	return queueClient
}
