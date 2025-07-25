package core

import (
	"context"

	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, prompt string) (domain.Task, error)
	ProcessTaskResult(ctx context.Context, event contracts.GenerateImageEvent) (domain.Task, error)
}

type TaskRepository interface {
	SaveTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTaskByID(ctx context.Context, id uint) (domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

type MessagePublisher interface {
	Publish(ctx context.Context, body []byte, queueName, correlationID string) error
}
