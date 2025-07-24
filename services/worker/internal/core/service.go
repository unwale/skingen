package core

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/unwale/skingen/pkg/contracts"
	cm "github.com/unwale/skingen/pkg/messaging"
	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
	"github.com/unwale/skingen/services/worker/internal/config"
)

type workerServiceImpl struct {
	modelServer ModelServer
	s3Client    S3Client
	publisher   *cm.RabbitMQPublisher
	config      *config.Config
}

func NewWorkerService(modelServer ModelServer, s3Client S3Client, publisher *cm.RabbitMQPublisher, cfg *config.Config) WorkerService {
	return &workerServiceImpl{
		modelServer: modelServer,
		s3Client:    s3Client,
		publisher:   publisher,
		config:      cfg,
	}
}

func (w *workerServiceImpl) GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error) {
	response, err := w.modelServer.GenerateImage(ctx, &pb.GenerateImageRequest{
		Prompt: request.Prompt,
	})
	if err != nil {
		return nil, err
	}

	objectId := uuid.New().String() + ".png"

	if err := w.s3Client.Upload(ctx, w.config.S3Config.Bucket, objectId, response.ImageData); err != nil {
		return nil, err
	}

	if err != nil {
		return &contracts.GenerateImageEvent{
			ObjectID: objectId,
			TaskID:   request.TaskID,
			Status:   contracts.TaskStatusFailed,
		}, err
	}

	event := &contracts.GenerateImageEvent{
		ObjectID: objectId,
		TaskID:   request.TaskID,
		Status:   contracts.TaskStatusCompleted,
	}
	eventBody, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	if err := w.publisher.Publish(ctx, eventBody, w.config.QueueConfig.TaskResultQueue); err != nil {
		return nil, err
	}

	return event, nil
}
