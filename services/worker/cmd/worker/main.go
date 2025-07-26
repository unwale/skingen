package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/minio/minio-go/v7"
	creds "github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/unwale/skingen/pkg/logging"
	cm "github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/worker/internal/adapters"
	"github.com/unwale/skingen/services/worker/internal/config"
	"github.com/unwale/skingen/services/worker/internal/core"
	"github.com/unwale/skingen/services/worker/internal/messaging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration")
	}

	logger := logging.NewLogger(cfg.ServiceName, cfg.LoggingLevel)
	logger.Info("Starting worker service", "port", cfg.Port)

	minioClient, err := minio.New(cfg.S3Config.Endpoint, &minio.Options{
		Creds:  creds.NewStaticV4(cfg.S3Config.AccessKey, cfg.S3Config.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Error("Failed to create MinIO client", "error", err)
		return
	}
	s3Client := adapters.NewS3ClientAdapter(minioClient)

	queueManager := cm.NewRabbitMQManager(cfg.RabbitMQUrl, logger)
	queueManager.Connect()
	defer queueManager.Close()

	conn, err := grpc.NewClient(cfg.ModelServerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to model server", "error", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Failed to close gRPC connection", "error", err)
		} else {
			logger.Info("gRPC connection closed successfully")
		}
	}()

	modelServerAdapter := adapters.NewTritonAdapter(conn)
	publisher := cm.NewRabbitMQPublisher(queueManager, logger)
	service := core.NewWorkerService(modelServerAdapter, s3Client, publisher, cfg, logger)
	taskCommandHandler := messaging.CreateTaskCommandHandler(service, logger)
	taskConsumer := cm.NewMessageConsumer(
		queueManager,
		cfg.QueueConfig.GenerationQueue,
		taskCommandHandler,
		logger,
	)
	if err := taskConsumer.Start(); err != nil {
		logger.Error("Failed to start message consumer", "error", err)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Info("Shutting down worker service gracefully")
	if err := taskConsumer.Shutdown(); err != nil {
		logger.Error("Failed to shutdown message consumer", "error", err)
	}
	queueManager.Close()
	logger.Info("Worker service shutdown complete")
}
