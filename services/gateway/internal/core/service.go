package core

import (
	"context"
	"log/slog"

	"github.com/unwale/skingen/pkg/contextutil"
	taskpb "github.com/unwale/skingen/services/task-service/generated/task/v1"
)

type gatewayServiceImpl struct {
	taskService TaskServicePort
	logger      *slog.Logger
}

func NewGatewayService(taskService TaskServicePort, logger *slog.Logger) GatewayService {
	return &gatewayServiceImpl{
		taskService: taskService,
		logger:      logger,
	}
}

func (s *gatewayServiceImpl) CreateTask(ctx context.Context, prompt string) (int, error) {
	logger := contextutil.FromContextOrDefault(ctx, s.logger)

	logger.Info("Calling Task Service", "method", "CreateTask", "prompt", prompt)

	req := &taskpb.CreateTaskRequest{
		Prompt: prompt,
	}

	resp, err := s.taskService.CreateTask(ctx, req)
	if err != nil {
		logger.Error("Task Service call failed", "error", err)
		return -1, err
	}

	logger.Info("Task created successfully", "task_id", resp.GetTaskId())
	return int(resp.GetTaskId()), nil
}
