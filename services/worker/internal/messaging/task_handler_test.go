package messaging

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/unwale/skingen/pkg/contracts"
)

type mockWorkerService struct {
	mock.Mock
}

func (m *mockWorkerService) GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.GenerateImageEvent), args.Error(1)
}

func TestCreateTaskCommandHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockService := new(mockWorkerService)
		command := &contracts.GenerateImageCommand{
			TaskID: 1,
			Prompt: "test prompt",
		}
		event := &contracts.GenerateImageEvent{
			ObjectID: "test-object-id",
			TaskID:   command.TaskID,
			Status:   contracts.TaskStatusCompleted,
		}

		mockService.On("GenerateImage", mock.Anything, command).Return(event, nil)

		handler := CreateTaskCommandHandler(mockService, logger)

		msgBody, _ := json.Marshal(command)
		msg := amqp091.Delivery{Body: msgBody}

		err := handler(msg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockService := new(mockWorkerService)
		command := &contracts.GenerateImageCommand{
			TaskID: 1,
			Prompt: "test prompt",
		}

		mockService.On("GenerateImage", mock.Anything, command).Return(nil, assert.AnError)

		handler := CreateTaskCommandHandler(mockService, logger)

		msgBody, _ := json.Marshal(command)
		msg := amqp091.Delivery{Body: msgBody}

		err := handler(msg)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		mockService.AssertExpectations(t)
	})
}
