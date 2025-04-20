package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"os"
	"testing"
	"time"

	"github.com/Egorpalan/workmate-test/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Setup()
	code := m.Run()
	os.Exit(code)
}

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)

	task.ID = "mock-id"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*entity.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) List(ctx context.Context, limit, offset int) ([]*entity.Task, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Task), args.Error(1)
}

// TestCreateTask тестирует создание задачи
func TestCreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockProcess := func(ctx context.Context) (json.RawMessage, error) {
		return json.RawMessage(`{"result":"success"}`), nil
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Task")).Return(nil)
	mockRepo.On("GetByID", mock.Anything, "mock-id").Return(&entity.Task{ID: "mock-id"}, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Task")).Return(nil)

	useCase := NewTaskUseCase(mockRepo, mockProcess)

	task, err := useCase.CreateTask(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, task)

	time.Sleep(100 * time.Millisecond)

	mockRepo.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*entity.Task"))
	mockRepo.AssertCalled(t, "GetByID", mock.Anything, "mock-id")
	mockRepo.AssertCalled(t, "Update", mock.Anything, mock.AnythingOfType("*entity.Task"))
}

// TestGetTaskByID тестирует получение задачи по ID
func TestGetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockProcess := func(ctx context.Context) (json.RawMessage, error) {
		return nil, nil
	}

	expectedTask := &entity.Task{ID: "test-id"}

	mockRepo.On("GetByID", mock.Anything, "test-id").Return(expectedTask, nil)

	useCase := NewTaskUseCase(mockRepo, mockProcess)

	task, err := useCase.GetTaskByID(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
	mockRepo.AssertExpectations(t)
}

// TestGetTaskByID_Error тестирует ошибку при получении задачи
func TestGetTaskByID_Error(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockProcess := func(ctx context.Context) (json.RawMessage, error) {
		return json.RawMessage(`{"result":"success"}`), nil
	}

	mockRepo.On("GetByID", mock.Anything, "task-id").Return(nil, errors.New("task not found"))

	useCase := NewTaskUseCase(mockRepo, mockProcess)

	task, err := useCase.GetTaskByID(context.Background(), "task-id")

	assert.Error(t, err)
	assert.Nil(t, task)

	mockRepo.AssertExpectations(t)
}

// TestListTasks тестирует получение списка задач
func TestListTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockProcess := func(ctx context.Context) (json.RawMessage, error) {
		return json.RawMessage(`{"result":"success"}`), nil
	}

	expectedTasks := []*entity.Task{
		{
			ID:        "task-1",
			Status:    entity.TaskStatusCompleted,
			Result:    json.RawMessage(`{"result":"success"}`),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "task-2",
			Status:    entity.TaskStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("List", mock.Anything, 10, 0).Return(expectedTasks, nil)

	useCase := NewTaskUseCase(mockRepo, mockProcess)

	tasks, err := useCase.ListTasks(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)

	mockRepo.AssertExpectations(t)
}
