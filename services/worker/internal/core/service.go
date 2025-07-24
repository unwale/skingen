package core

import (
	"context"
	"encoding/json"

	"github.com/unwale/skingen/pkg/contracts"
	cm "github.com/unwale/skingen/pkg/messaging"
	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
	"github.com/unwale/skingen/services/worker/internal/config"
)

type workerServiceImpl struct {
	modelServer ModelServer
	publisher   *cm.RabbitMQPublisher
	queueConfig config.QueueConfig
}

func NewWorkerService(modelServer ModelServer, publisher *cm.RabbitMQPublisher, cfg config.QueueConfig) WorkerService {
	return &workerServiceImpl{
		modelServer: modelServer,
		publisher:   publisher,
		queueConfig: cfg,
	}
}

func (w *workerServiceImpl) GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error) {
	response, err := w.modelServer.GenerateImage(ctx, &pb.GenerateImageRequest{
		Prompt: request.Prompt,
	})
	if err != nil {
		return &contracts.GenerateImageEvent{
			ImageURL: "",
			TaskID:   request.TaskID,
			Status:   contracts.TaskStatusFailed,
		}, err
	}

	event := &contracts.GenerateImageEvent{
		ImageURL: string(response.ImageData),
		TaskID:   request.TaskID,
		Status:   contracts.TaskStatusCompleted,
	}
	eventBody, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	w.publisher.Publish(ctx, eventBody, w.queueConfig.TaskResultQueue)

	return event, nil
}
