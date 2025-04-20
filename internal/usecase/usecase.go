package usecase

import (
	"context"

	"github.com/Egorpalan/workmate-test/internal/entity"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context) (*entity.Task, error)
	GetTaskByID(ctx context.Context, id string) (*entity.Task, error)
	ListTasks(ctx context.Context, limit, offset int) ([]*entity.Task, error)
}

type UseCase struct {
	Task TaskUseCase
}

// NewUseCase создает новый экземпляр UseCase
func NewUseCase(task TaskUseCase) *UseCase {
	return &UseCase{
		Task: task,
	}
}
