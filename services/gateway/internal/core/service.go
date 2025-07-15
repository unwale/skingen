package core

import (
	"context"

	taskpb "github.com/unwale/skingen/services/task-service/generated/task/v1"
)

type gatewayServiceImpl struct {
	taskService TaskServicePort
}

func NewGatewayService(taskService TaskServicePort) GatewayService {
	return &gatewayServiceImpl{
		taskService: taskService,
	}
}

func (s *gatewayServiceImpl) CreateTask(ctx context.Context, prompt string) (int, error) {
	req := &taskpb.CreateTaskRequest{
		Prompt: prompt,
	}

	resp, err := s.taskService.CreateTask(ctx, req)
	if err != nil {
		return -1, err
	}

	return int(resp.GetTaskId()), nil
}
