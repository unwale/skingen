package core

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"github.com/unwale/skingen/pkg/contextutil"
	"github.com/unwale/skingen/pkg/contracts"
	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
	"github.com/unwale/skingen/services/worker/internal/config"
)

type workerServiceImpl struct {
	modelServer ModelServer
	s3Client    S3Client
	publisher   MessagePublisher
	config      *config.Config
	logger      *slog.Logger
}

func NewWorkerService(modelServer ModelServer, s3Client S3Client, publisher MessagePublisher, cfg *config.Config, logger *slog.Logger) WorkerService {
	return &workerServiceImpl{
		modelServer: modelServer,
		s3Client:    s3Client,
		publisher:   publisher,
		config:      cfg,
		logger:      logger,
	}
}

func (w *workerServiceImpl) GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error) {
	logger := contextutil.FromContextOrDefault(ctx, w.logger)

	response, err := w.modelServer.GenerateImage(ctx, &pb.GenerateImageRequest{
		Prompt: request.Prompt,
	})
	if err != nil {
		logger.Error("Failed to generate image", "error", err)
		return nil, err
	}

	objectId := uuid.New().String() + ".png"

	logger.Info("Uploading image to S3", "object_id", objectId)
	if err := w.s3Client.Upload(ctx, w.config.S3Config.Bucket, objectId, response.ImageData); err != nil {
		logger.Error("Failed to upload image to S3", "error", err, "object_id", objectId)
		return nil, err
	}

	event := &contracts.GenerateImageEvent{
		ObjectID: objectId,
		TaskID:   request.TaskID,
		Status:   contracts.TaskStatusCompleted,
	}
	eventBody, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal event", "error", err)
		return nil, err
	}

	logger.Info("Publishing image generation event", "task_id", event.TaskID)
	if err := w.publisher.Publish(
		ctx,
		eventBody,
		w.config.QueueConfig.TaskResultQueue,
		contextutil.CorrelationIDFromContext(ctx),
	); err != nil {
		logger.Error("Failed to publish event", "error", err)
		return nil, err
	}

	return event, nil
}
