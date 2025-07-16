package core

import (
	"context"

	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type taskServiceImpl struct {
	repo      TaskRepository
	publisher MessagePublisher
}

func NewTaskService(repo TaskRepository, publisher MessagePublisher) TaskService {
	return &taskServiceImpl{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, prompt string) (domain.Task, error) {
	task := domain.Task{Prompt: prompt}
	savedTask, err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	return savedTask, nil
}
