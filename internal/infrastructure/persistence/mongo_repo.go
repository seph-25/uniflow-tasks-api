package persistence

import (
	"context"
	"fmt"
	"time"
	"uniflow-api/internal/application/ports"
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

// DueToday obtiene todas las tareas con fecha vencimiento hoy
func (r *MongoTaskRepository) DueToday(ctx context.Context, userID string, loc *time.Location) ([]domain.Task, error) {
	// Calcular inicio y fin del día en la zona horaria especificada
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"userId": userID,
		"dueDate": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	opts := options.Find()
	opts.SetSort(bson.M{"dueDate": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error al buscar tareas de hoy: %w", err)
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

// Find obtiene tareas con filtros complejos y paginación
func (r *MongoTaskRepository) Find(ctx context.Context, f ports.TaskFilter) ([]domain.Task, ports.PageInfo, error) {
	// Construir filtro base
	filter := bson.M{
		"userId": f.UserID,
	}

	// Filtros opcionales
	if f.Status != "" {
		if !isValidStatus(f.Status) {
			return nil, ports.PageInfo{}, fmt.Errorf("estado inválido: %s", f.Status)
		}
		filter["status"] = f.Status
	}

	if f.SubjectID != "" {
		filter["subjectId"] = f.SubjectID
	}

	if f.PeriodID != "" {
		filter["periodId"] = f.PeriodID
	}

	// Rango de fechas
	dateFilter := bson.M{}
	if f.From != nil {
		dateFilter["$gte"] = f.From
	}
	if f.To != nil {
		dateFilter["$lte"] = f.To
	}
	if len(dateFilter) > 0 {
		filter["dueDate"] = dateFilter
	}

	// Búsqueda de texto (si existe query)
	if f.Query != "" {
		filter["$or"] = bson.A{
			bson.M{"title": bson.M{"$regex": f.Query, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": f.Query, "$options": "i"}},
		}
	}

	// Contar total de documentos
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, ports.PageInfo{}, fmt.Errorf("error al contar tareas: %w", err)
	}

	// Opciones de paginación y ordenamiento
	opts := options.Find()
	opts.SetSkip(int64(f.Offset))
	opts.SetLimit(int64(f.Limit))
	opts.SetSort(bson.M{"dueDate": 1})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, ports.PageInfo{}, fmt.Errorf("error al buscar tareas: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, ports.PageInfo{}, fmt.Errorf("error al decodificar tareas: %w", err)
	}

	if tasks == nil {
		tasks = []domain.Task{}
	}

	// Calcular metadata de paginación
	totalPages := int(total) / f.Limit
	if int(total)%f.Limit > 0 {
		totalPages++
	}

	page := f.Offset/f.Limit + 1
	pageInfo := ports.PageInfo{
		Total:      total,
		Page:       page,
		Limit:      f.Limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return tasks, pageInfo, nil
}

// Search busca tareas por texto libre (reutiliza Find)
func (r *MongoTaskRepository) Search(ctx context.Context, f ports.TaskFilter) ([]domain.Task, ports.PageInfo, error) {
	// Search es básicamente Find con énfasis en el Query
	// Si no hay Query, retornar error
	if f.Query == "" {
		return nil, ports.PageInfo{}, fmt.Errorf("query de búsqueda vacía")
	}

	return r.Find(ctx, f)
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
