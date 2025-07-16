package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error) {
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
	taskService := NewTaskService(mockRepo, mockPublisher)

	prompt := "test prompt"
	expectedTask := domain.Task{ID: 1, Prompt: prompt}

	mockRepo.On("SaveTask", mock.Anything, mock.AnythingOfType("domain.Task")).Return(expectedTask, nil)

	task, err := taskService.CreateTask(context.Background(), prompt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.ID != expectedTask.ID || task.Prompt != expectedTask.Prompt {
		t.Errorf("expected task %v, got %v", expectedTask, task)
	}

	mockRepo.AssertExpectations(t)
}
