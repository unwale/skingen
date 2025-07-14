package core

import (
	"context"

	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, prompt string) (domain.Task, error)
}

type TaskRepository interface {
	SaveTask(ctx context.Context, task domain.Task) (domain.Task, error)
}
