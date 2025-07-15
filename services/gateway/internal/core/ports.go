package core

import (
	"context"

	taskpb "github.com/unwale/skingen/services/task-service/generated/task/v1"
)

type GatewayService interface {
	CreateTask(ctx context.Context, prompt string) (int, error)
}

type TaskServicePort interface {
	CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error)
}
