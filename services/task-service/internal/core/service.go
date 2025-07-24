package core

import (
	"context"
	"encoding/json"

	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/domain"
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
	task := domain.Task{
		Prompt: prompt,
		Status: domain.TaskStatusPending,
	}
	savedTask, err := s.repo.SaveTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	generateImageCommand := contracts.GenerateImageCommand{
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

func (s *taskServiceImpl) ProcessTaskResult(ctx context.Context, event contracts.GenerateImageEvent) (domain.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, event.TaskID)
	if err != nil {
		return domain.Task{}, err
	}

	if event.Status != contracts.TaskStatusCompleted {
		task.Status = domain.TaskStatusFailed
	} else {
		task.Status = domain.TaskStatusCompleted
	}
	task.ObjectID = event.ObjectID
	updatedTask, err := s.repo.UpdateTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}

	return updatedTask, nil
}
