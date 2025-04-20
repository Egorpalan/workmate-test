package postgresql

import (
	"context"
	"fmt"

	"github.com/Egorpalan/workmate-test/internal/entity"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type TaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository создает новый экземпляр TaskRepository
func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Create создает новую задачу в базе данных
func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	query := `
        INSERT INTO tasks (status, result, error)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `

	row := r.db.QueryRowxContext(
		ctx,
		query,
		task.Status,
		task.Result,
		task.Error,
	)

	err := row.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetByID возвращает задачу по ее ID
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	query := `
        SELECT id, status, result, error, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	var task entity.Task
	err := r.db.GetContext(ctx, &task, query, id)
	if err != nil {
		logger.Error("Failed to get task by ID", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &task, nil
}

// Update обновляет существующую задачу
func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	query := `
        UPDATE tasks
        SET status = $1, result = $2, error = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING updated_at
    `

	row := r.db.QueryRowContext(
		ctx,
		query,
		task.Status,
		task.Result,
		task.Error,
		task.ID,
	)

	err := row.Scan(&task.UpdatedAt)
	if err != nil {
		logger.Error("Failed to update task", zap.String("id", task.ID), zap.Error(err))
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// List возвращает список задач с пагинацией
func (r *TaskRepository) List(ctx context.Context, limit, offset int) ([]*entity.Task, error) {
	query := `
        SELECT id, status, result, error, created_at, updated_at
        FROM tasks
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	tasks := make([]*entity.Task, 0)
	err := r.db.SelectContext(ctx, &tasks, query, limit, offset)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}
