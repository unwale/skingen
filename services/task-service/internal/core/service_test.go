package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockRepository) GetTaskByID(ctx context.Context, id uint) (domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockRepository) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

type mockPublisher struct {
	mock.Mock
}

func (m *mockPublisher) Publish(ctx context.Context, body []byte, queueName string) error {
	args := m.Called(ctx, body, queueName)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	mockRepo := new(mockRepository)
	mockPublisher := new(mockPublisher)
	queueConfig := config.QueueConfig{
		GenerateImageQueue: "generate_image_queue",
		TaskResultQueue:    "task_result_queue",
	}
	taskService := NewTaskService(mockRepo, mockPublisher, queueConfig)

	prompt := "test prompt"
	expectedTask := domain.Task{ID: 1, Prompt: prompt}

	mockRepo.On("SaveTask", mock.Anything, mock.AnythingOfType("domain.Task")).Return(expectedTask, nil)
	mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("[]uint8"), queueConfig.GenerateImageQueue).Return(nil)

	task, err := taskService.CreateTask(context.Background(), prompt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != expectedTask.ID || task.Prompt != expectedTask.Prompt {
		t.Errorf("expected task %v, got %v", expectedTask, task)
	}

	mockRepo.AssertExpectations(t)
}
