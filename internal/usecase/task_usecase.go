package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Egorpalan/workmate-test/internal/entity"
	"github.com/Egorpalan/workmate-test/internal/repository"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"go.uber.org/zap"
)

// LongRunningTask представляет функцию, выполняющую длительную задачу
type LongRunningTask func(ctx context.Context) (json.RawMessage, error)

type taskUseCase struct {
	taskRepo    repository.TaskRepository
	processTask LongRunningTask
}

// NewTaskUseCase создает новый экземпляр taskUseCase
func NewTaskUseCase(taskRepo repository.TaskRepository, processTask LongRunningTask) *taskUseCase {
	return &taskUseCase{
		taskRepo:    taskRepo,
		processTask: processTask,
	}
}

// CreateTask создает новую задачу и запускает ее асинхронное выполнение
func (u *taskUseCase) CreateTask(ctx context.Context) (*entity.Task, error) {
	task := &entity.Task{
		Status: entity.TaskStatusPending,
		Result: json.RawMessage([]byte("{}")), // Пустой JSON
	}

	if err := u.taskRepo.Create(ctx, task); err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	go u.processTaskAsync(task.ID)

	return task, nil
}

// GetTaskByID возвращает задачу по ее ID
func (u *taskUseCase) GetTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	task, err := u.taskRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get task by ID", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return task, nil
}

// ListTasks возвращает список задач с пагинацией
func (u *taskUseCase) ListTasks(ctx context.Context, limit, offset int) ([]*entity.Task, error) {
	tasks, err := u.taskRepo.List(ctx, limit, offset)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

// processTaskAsync выполняет задачу асинхронно
func (u *taskUseCase) processTaskAsync(taskID string) {
	ctx := context.Background()

	task, err := u.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		logger.Error("Failed to get task by ID for processing", zap.String("id", taskID), zap.Error(err))
		return
	}

	task.Status = entity.TaskStatusProcessing
	if err := u.taskRepo.Update(ctx, task); err != nil {
		logger.Error("Failed to update task status to processing", zap.String("id", taskID), zap.Error(err))
		return
	}

	result, err := u.processTask(ctx)

	task.UpdatedAt = time.Now()
	if err != nil {
		task.Status = entity.TaskStatusFailed
		task.Error = err.Error()
	} else {
		task.Status = entity.TaskStatusCompleted
		task.Result = result
	}

	if err := u.taskRepo.Update(ctx, task); err != nil {
		logger.Error("Failed to update task with result", zap.String("id", taskID), zap.Error(err))
	}
}
