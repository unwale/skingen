package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/pkg/contracts"
	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) CreateTask(ctx context.Context, prompt string) (domain.Task, error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskService) ProcessTaskResult(ctx context.Context, event contracts.GenerateImageEvent) (domain.Task, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(domain.Task), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockTaskService := new(mockTaskService)
		handler := NewHandler(mockTaskService)

		prompt := "test prompt"
		expectedTask := domain.Task{ID: 1, Prompt: prompt}
		mockTaskService.On("CreateTask", mock.Anything, prompt).Return(expectedTask, nil)

		req := &pb.CreateTaskRequest{Prompt: prompt}
		resp, err := handler.CreateTask(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, uint32(expectedTask.ID), resp.TaskId)
	})

	t.Run("error", func(t *testing.T) {
		mockTaskService := new(mockTaskService)
		handler := NewHandler(mockTaskService)

		prompt := "test prompt"
		mockTaskService.On("CreateTask", mock.Anything, prompt).Return(domain.Task{}, assert.AnError)

		req := &pb.CreateTaskRequest{Prompt: prompt}
		resp, err := handler.CreateTask(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
