package core

import (
	"context"
	"encoding/json"

	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/domain"
	"github.com/unwale/skingen/services/task-service/internal/messaging"
)

type taskServiceImpl struct {
	repo        TaskRepository
	publisher   MessagePublisher
	queueConfig config.QueueConfig
}

func NewTaskService(repo TaskRepository, publisher MessagePublisher, cfg config.QueueConfig) TaskService {
	return &taskServiceImpl{
		repo:        repo,
		publisher:   publisher,
		queueConfig: cfg,
	}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, prompt string) (domain.Task, error) {
	task := domain.Task{Prompt: prompt}
	savedTask, err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	generateImageCommand := messaging.GenerateImageCommand{
		TaskID: savedTask.ID,
		Prompt: savedTask.Prompt,
	}

	body, err := json.Marshal(generateImageCommand)
	if err != nil {
		return domain.Task{}, err
	}
	err = s.publisher.Publish(ctx, body, s.queueConfig.GenerateImageQueue)
	if err != nil {
		return domain.Task{}, err
	}

	return savedTask, nil
}
