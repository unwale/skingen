package core

import (
	"context"
	"io"
	"log/slog"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/unwale/skingen/pkg/contracts"
	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
	"github.com/unwale/skingen/services/worker/internal/config"
)

type mockModelServer struct {
	mock.Mock
}

func (m *mockModelServer) GenerateImage(ctx context.Context, request *pb.GenerateImageRequest) (*pb.GenerateImageResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GenerateImageResponse), args.Error(1)
}

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) Upload(ctx context.Context, bucket, key string, data []byte) error {
	args := m.Called(ctx, bucket, key, data)
	return args.Error(0)
}

type mockRabbitMQPublisher struct {
	mock.Mock
}

func (m *mockRabbitMQPublisher) Publish(ctx context.Context, body []byte, queueName, correlationID string) error {
	args := m.Called(ctx, body, queueName)
	return args.Error(0)
}

func TestGenerateImage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockModelServer := new(mockModelServer)
		mockS3Client := new(mockS3Client)
		mockPublisher := new(mockRabbitMQPublisher)
		cfg := &config.Config{
			S3Config: config.S3Config{
				Bucket: "test-bucket",
			},
			QueueConfig: config.QueueConfig{
				TaskResultQueue: "task-result-queue",
			},
		}

		service := NewWorkerService(mockModelServer, mockS3Client, mockPublisher, cfg, logger)

		request := &contracts.GenerateImageCommand{
			Prompt: "test prompt",
			TaskID: 123,
		}

		mockModelServer.On("GenerateImage", mock.Anything, mock.Anything).Return(&pb.GenerateImageResponse{ImageData: []byte("image data")}, nil)
		mockS3Client.On("Upload", mock.Anything, cfg.S3Config.Bucket, mock.AnythingOfType("string"), []byte("image data")).Return(nil)
		mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("[]uint8"), cfg.QueueConfig.TaskResultQueue, mock.Anything).Return(nil)

		event, err := service.GenerateImage(context.Background(), request)

		require.NoError(t, err)
		assert.NotNil(t, event)
		assert.Equal(t, contracts.TaskStatusCompleted, event.Status)
		assert.Contains(t, event.ObjectID, ".png")

		mockModelServer.AssertExpectations(t)
		mockS3Client.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("model server error", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockModelServer := new(mockModelServer)
		mockS3Client := new(mockS3Client)
		mockPublisher := new(mockRabbitMQPublisher)
		cfg := &config.Config{}

		service := NewWorkerService(mockModelServer, mockS3Client, mockPublisher, cfg, logger)

		request := &contracts.GenerateImageCommand{
			Prompt: "test prompt",
			TaskID: 123,
		}

		mockModelServer.On("GenerateImage", mock.Anything, mock.Anything).Return(nil, assert.AnError)

		event, err := service.GenerateImage(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, event)

		mockModelServer.AssertExpectations(t)
		mockS3Client.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("s3 upload error", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockModelServer := new(mockModelServer)
		mockS3Client := new(mockS3Client)
		mockPublisher := new(mockRabbitMQPublisher)
		cfg := &config.Config{
			S3Config: config.S3Config{
				Bucket: "test-bucket",
			},
			QueueConfig: config.QueueConfig{
				TaskResultQueue: "task-result-queue",
			},
		}

		service := NewWorkerService(mockModelServer, mockS3Client, mockPublisher, cfg, logger)

		request := &contracts.GenerateImageCommand{
			Prompt: "test prompt",
			TaskID: 123,
		}

		mockModelServer.On("GenerateImage", mock.Anything, mock.Anything).Return(&pb.GenerateImageResponse{ImageData: []byte("image data")}, nil)
		mockS3Client.On("Upload", mock.Anything, cfg.S3Config.Bucket, mock.AnythingOfType("string"), []byte("image data")).Return(assert.AnError)

		event, err := service.GenerateImage(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, event)

		mockModelServer.AssertExpectations(t)
		mockS3Client.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("publisher error", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockModelServer := new(mockModelServer)
		mockS3Client := new(mockS3Client)
		mockPublisher := new(mockRabbitMQPublisher)
		cfg := &config.Config{
			S3Config: config.S3Config{
				Bucket: "test-bucket",
			},
			QueueConfig: config.QueueConfig{
				TaskResultQueue: "task-result-queue",
			},
		}

		service := NewWorkerService(mockModelServer, mockS3Client, mockPublisher, cfg, logger)

		request := &contracts.GenerateImageCommand{
			Prompt: "test prompt",
			TaskID: 123,
		}

		mockModelServer.On("GenerateImage", mock.Anything, mock.Anything).Return(&pb.GenerateImageResponse{ImageData: []byte("image data")}, nil)
		mockS3Client.On("Upload", mock.Anything, cfg.S3Config.Bucket, mock.AnythingOfType("string"), []byte("image data")).Return(nil)
		mockPublisher.On("Publish", mock.Anything, mock.AnythingOfType("[]uint8"), cfg.QueueConfig.TaskResultQueue, mock.Anything).Return(assert.AnError)

		event, err := service.GenerateImage(context.Background(), request)

		assert.Error(t, err)
		assert.Nil(t, event)

		mockModelServer.AssertExpectations(t)
		mockS3Client.AssertExpectations(t)

		mockPublisher.AssertExpectations(t)
	})
}
