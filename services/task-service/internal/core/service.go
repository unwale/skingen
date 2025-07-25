package core

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/unwale/skingen/pkg/contextutil"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type taskServiceImpl struct {
	repo        TaskRepository
	publisher   MessagePublisher
	queueConfig config.QueueConfig
	logger      *slog.Logger
}

func NewTaskService(repo TaskRepository, publisher MessagePublisher, cfg config.QueueConfig, logger *slog.Logger) TaskService {
	return &taskServiceImpl{
		repo:        repo,
		publisher:   publisher,
		queueConfig: cfg,
		logger:      logger,
	}
}

func (s *taskServiceImpl) CreateTask(ctx context.Context, prompt string) (domain.Task, error) {
	logger := contextutil.FromContextOrDefault(ctx, s.logger)

	logger.Info("Creating task")

	task := domain.Task{
		Prompt: prompt,
		Status: domain.TaskStatusPending,
	}
	savedTask, err := s.repo.SaveTask(ctx, task)
	if err != nil {
		logger.Error("Failed to save task", "error", err)
		return domain.Task{}, err
	}

	logger.Info("Task saved successfully", "task_id", savedTask.ID)

	generateImageCommand := contracts.GenerateImageCommand{
		TaskID: savedTask.ID,
		Prompt: savedTask.Prompt,
	}

	body, err := json.Marshal(generateImageCommand)
	if err != nil {
		logger.Error("Failed to marshal generate image command", "error", err)
		return domain.Task{}, err
	}
	err = s.publisher.Publish(ctx, body, s.queueConfig.GenerateImageQueue, contextutil.CorrelationIDFromContext(ctx))
	if err != nil {
		logger.Error("Failed to publish message to queue", "error", err)
		return domain.Task{}, err
	}

	logger.Info("Published generate image command to queue", "task_id", savedTask.ID)

	return savedTask, nil
}

func (s *taskServiceImpl) ProcessTaskResult(ctx context.Context, event contracts.GenerateImageEvent) (domain.Task, error) {
	logger := contextutil.FromContextOrDefault(ctx, s.logger)

	task, err := s.repo.GetTaskByID(ctx, event.TaskID)
	if err != nil {
		logger.Error("Failed to get task by ID", "task_id", event.TaskID, "error", err)
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
		logger.Error("Failed to update task", "task_id", task.ID, "error", err)
		return domain.Task{}, err
	}

	logger.Info("Task updated successfully", "task_id", updatedTask.ID, "status", updatedTask.Status)

	return updatedTask, nil
}
