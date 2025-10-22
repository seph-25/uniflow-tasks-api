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
	"uniflow-api/internal/infrastructure/handlers"
	"uniflow-api/internal/infrastructure/persistence"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}
	// Configurar puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configurar modo Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Obtener MONGO_URI del ambiente
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Println("ERROR: MONGO_URI no configurada en variables de entorno")
		os.Exit(1)
	}

	// Obtener nombre de base de datos del ambiente
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "my-culster-name-ds2025" // Default si no se proporciona
	}

	// Conectar a MongoDB con contexto y timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("ERROR al conectar a MongoDB: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.Background())

	// Verificar conexión (ping)
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Printf("ERROR al hacer ping a MongoDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Conectado a MongoDB exitosamente")

	// Seleccionar base de datos y colección
	database := client.Database(mongoDB)
	tasksCollection := database.Collection("tasks")

	// Crear repositorio MongoDB
	taskRepo := persistence.NewMongoTaskRepository(tasksCollection)

	// Crear servicio CON inyección del repositorio
	taskService := application.NewTaskService(taskRepo)

	// Crear router
	r := gin.Default()

	// Crear handlers
	taskHandler := handlers.NewTaskHandler(taskService)

	// Rutas
	r.GET("/health", handlers.HealthHandler)
	r.GET("/tasks", taskHandler.GetTasks)

	// Levantar servidor
	fmt.Printf("Servidor escuchando en puerto %s\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("ERROR al levantar servidor: %v\n", err)
		os.Exit(1)
	}
}