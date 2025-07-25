package core

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	taskpb "github.com/unwale/skingen/services/task-service/generated/task/v1"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*taskpb.CreateTaskResponse), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	mockTaskService := new(mockTaskService)
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	gatewayService := NewGatewayService(mockTaskService, testLogger)

	prompt := "test prompt"
	req := &taskpb.CreateTaskRequest{Prompt: prompt}
	resp := &taskpb.CreateTaskResponse{TaskId: 1}

	mockTaskService.On("CreateTask", mock.Anything, req).Return(resp, nil)

	taskID, err := gatewayService.CreateTask(context.Background(), prompt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if taskID != int(resp.GetTaskId()) {
		t.Errorf("expected task ID %d, got %d", resp.GetTaskId(), taskID)
	}

	mockTaskService.AssertExpectations(t)
	mockTaskService.AssertCalled(t, "CreateTask", mock.Anything, req)
	mockTaskService.AssertNotCalled(t, "CreateTask", mock.Anything, &taskpb.CreateTaskRequest{Prompt: "different prompt"})
}
