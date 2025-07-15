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

func TestCreateTask(t *testing.T) {
	mockRepo := new(mockRepository)
	taskService := NewTaskService(mockRepo)

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
