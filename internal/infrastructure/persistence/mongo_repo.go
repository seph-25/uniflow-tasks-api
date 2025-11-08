package persistence

import (
	"context"
	"fmt"
	"time"
	"uniflow-api/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoTaskRepository implementa TaskRepository usando MongoDB
type MongoTaskRepository struct {
	collection *mongo.Collection
}

// NewMongoTaskRepository crea una nueva instancia de MongoTaskRepository
func NewMongoTaskRepository(collection *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: collection,
	}
}

// Create inserta una nueva tarea en MongoDB
func (r *MongoTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	// Generar ObjectID si no existe
	if task.ID == "" {
		task.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("error al crear tarea: %w", err)
	}

	return nil
}

// GetByID obtiene una tarea específica por ID y verifica pertenencia al usuario
func (r *MongoTaskRepository) GetByID(ctx context.Context, taskID, userID string) (*domain.Task, error) {
	// Filtro: tarea debe pertenecer al usuario Y tener ese ID (como string)
	filter := bson.M{
		"_id":    taskID,
		"userId": userID,
	}

	var task domain.Task
	err := r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("tarea no encontrada")
		}
		return nil, fmt.Errorf("error al obtener tarea: %w", err)
	}

	return &task, nil
}

// GetAll obtiene todas las tareas de un usuario
func (r *MongoTaskRepository) GetAll(ctx context.Context, userID string) ([]domain.Task, error) {
	// Filtro: solo tareas del usuario actual
	filter := bson.M{
		"userId": userID,
	}
	// Opciones: ordenar por dueDate ascendente
	opts := options.Find()
	opts.SetSort(bson.M{"dueDate": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error al buscar tareas: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar tareas: %w", err)
	}

	// Si no hay tareas, retornar slice vacío (no nil)
	if tasks == nil {
		tasks = []domain.Task{}
	}

	return tasks, nil
}

// GetByUserAndStatus obtiene tareas filtradas por usuario y estado
func (r *MongoTaskRepository) GetByUserAndStatus(ctx context.Context, userID, status string) ([]domain.Task, error) {
	// Validar que el status sea válido
	if !isValidStatus(status) {
		return nil, fmt.Errorf("estado inválido: %s", status)
	}

	filter := bson.M{
		"userId": userID,
		"status": status,
	}

	opts := options.Find()
	opts.SetSort(bson.M{"dueDate": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error al buscar tareas por estado: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar tareas: %w", err)
	}

	if tasks == nil {
		tasks = []domain.Task{}
	}

	return tasks, nil
}

// Update actualiza una tarea existente
func (r *MongoTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	// Validar que la tarea puede ser modificada
	if err := task.CanBeModified(); err != nil {
		return fmt.Errorf("no se puede actualizar: %w", err)
	}

	// Filtro: asegurarse que pertenece al usuario (usando string ID)
	filter := bson.M{
		"_id":    task.ID,
		"userId": task.UserID,
	}

	// Update: reemplazar el documento
	result, err := r.collection.ReplaceOne(ctx, filter, task)
	if err != nil {
		return fmt.Errorf("error al actualizar tarea: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("tarea no encontrada o no pertenece al usuario")
	}

	return nil
}

// Delete elimina una tarea
func (r *MongoTaskRepository) Delete(ctx context.Context, taskID, userID string) error {
	// Validar que no esté completada (regla de negocio)
	task, err := r.GetByID(ctx, taskID, userID)
	if err != nil {
		return err
	}

	if err := task.CanBeDeleted(); err != nil {
		return err
	}

	filter := bson.M{
		"_id":    taskID,
		"userId": userID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error al eliminar tarea: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tarea no encontrada")
	}

	return nil
}

// Métodos para Fase 3 (stubs por ahora)
func (r *MongoTaskRepository) Find(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	return []domain.Task{}, domain.PageInfo{}, nil
}

func (r *MongoTaskRepository) DueToday(ctx context.Context, userID string, loc *time.Location) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (r *MongoTaskRepository) Search(ctx context.Context, f domain.TaskFilter) ([]domain.Task, domain.PageInfo, error) {
	return []domain.Task{}, domain.PageInfo{}, nil
}

func (r *MongoTaskRepository) Aggregated(ctx context.Context, userID string, until time.Time) (domain.Stats, error) {
	return domain.Stats{}, nil
}

// Helper: validar status
func isValidStatus(status string) bool {
	validStatuses := []string{
		domain.StatusTodo,
		domain.StatusInProgress,
		domain.StatusInReview,
		domain.StatusDone,
		domain.StatusCancelled,
	}
	for _, v := range validStatuses {
		if v == status {
			return true
		}
	}
	return false
}