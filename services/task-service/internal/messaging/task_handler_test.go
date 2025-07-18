package messaging

import (
	"context"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) ProcessTaskResult(ctx context.Context, event contracts.GenerateImageEvent) (domain.Task, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskService) CreateTask(ctx context.Context, prompt string) (domain.Task, error) {
	args := m.Called(ctx, prompt)
	return args.Get(0).(domain.Task), args.Error(1)
}

func TestCreateTaskResultHandler(t *testing.T) {
	service := new(mockTaskService)
	handler := CreateTaskResultHandler(service)

	event := contracts.GenerateImageEvent{
		TaskID:   1,
		Status:   domain.TaskStatusCompleted,
		ImageURL: "http://example.com/image.png",
	}
	msg := amqp091.Delivery{
		Body: []byte(`{"task_id":1,"status":"completed","image_url":"http://example.com/image.png"}`),
	}
	service.On("ProcessTaskResult", mock.Anything, event).Return(domain.Task{ID: 1}, nil)
	err := handler(msg)
	assert.NoError(t, err)
	service.AssertExpectations(t)
}
