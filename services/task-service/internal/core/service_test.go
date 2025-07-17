package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/pkg/contracts"
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

func TestProcessTaskResult(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockRepository)
		taskService := NewTaskService(mockRepo, nil, config.QueueConfig{})

		event := contracts.GenerateImageEvent{
			TaskID:   1,
			Status:   domain.TaskStatusCompleted,
			ImageURL: "image.png",
		}
		expectedTask := domain.Task{ID: 1, Status: domain.TaskStatusCompleted, ResultURL: event.ImageURL}

		mockRepo.On("GetTaskByID", mock.Anything, event.TaskID).Return(expectedTask, nil)
		mockRepo.On("UpdateTask", mock.Anything, expectedTask).Return(expectedTask, nil)

		task, err := taskService.ProcessTaskResult(context.Background(), event)
		assert.NoError(t, err)

		assert.Equal(t, expectedTask.ID, task.ID)
		assert.Equal(t, expectedTask.Status, task.Status)
		assert.Equal(t, expectedTask.ResultURL, task.ResultURL)
		assert.Equal(t, expectedTask.Prompt, task.Prompt)

		mockRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		mockRepo := new(mockRepository)
		taskService := NewTaskService(mockRepo, nil, config.QueueConfig{})

		event := contracts.GenerateImageEvent{
			TaskID: 1,
			Status: domain.TaskStatusCompleted,
		}

		mockRepo.On("GetTaskByID", mock.Anything, event.TaskID).Return(domain.Task{}, assert.AnError)

		task, err := taskService.ProcessTaskResult(context.Background(), event)
		assert.Error(t, err)
		assert.Equal(t, domain.Task{}, task)

		mockRepo.AssertExpectations(t)
	})

	t.Run("update task error", func(t *testing.T) {
		mockRepo := new(mockRepository)
		taskService := NewTaskService(mockRepo, nil, config.QueueConfig{})

		event := contracts.GenerateImageEvent{
			TaskID: 1,
			Status: domain.TaskStatusCompleted,
		}
		task := domain.Task{ID: 1, Status: domain.TaskStatusPending}
		taskCompleted := domain.Task{ID: 1, Status: domain.TaskStatusCompleted}

		mockRepo.On("GetTaskByID", mock.Anything, event.TaskID).Return(task, nil)
		mockRepo.On("UpdateTask", mock.Anything, taskCompleted).Return(domain.Task{}, assert.AnError)

		result, err := taskService.ProcessTaskResult(context.Background(), event)
		assert.Error(t, err)
		assert.Equal(t, domain.Task{}, result)

		mockRepo.AssertExpectations(t)
	})

}
