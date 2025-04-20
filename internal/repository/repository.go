package repository

import (
	"context"

	"github.com/Egorpalan/workmate-test/internal/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id string) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	List(ctx context.Context, limit, offset int) ([]*entity.Task, error)
}

type Repository struct {
	Task TaskRepository
}

// NewRepository создает новый экземпляр всех репозиториев
func NewRepository(task TaskRepository) *Repository {
	return &Repository{
		Task: task,
	}
}
