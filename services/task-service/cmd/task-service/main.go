package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	grpc_server "google.golang.org/grpc"

	"github.com/unwale/skingen/pkg/logging"
	cm "github.com/unwale/skingen/pkg/messaging"
	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"github.com/unwale/skingen/services/task-service/internal/api/grpc"
	"github.com/unwale/skingen/services/task-service/internal/api/grpc/interceptors"
	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/core"
	"github.com/unwale/skingen/services/task-service/internal/database"
	"github.com/unwale/skingen/services/task-service/internal/messaging"
	"github.com/unwale/skingen/services/task-service/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logging.NewLogger(cfg.ServiceName, cfg.LoggingLevel)
	slog.SetDefault(logger)
	logger.Info("Starting task service", "port", cfg.Port)

	queueManager := cm.NewRabbitMQManager(cfg.RabbitMQUrl, logger)
	queueManager.Connect()
	defer queueManager.Close()

	db, err := database.NewConnection(*cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}

	queuePublisher := cm.NewRabbitMQPublisher(queueManager, logger)
	repo := repository.NewTaskRepository(db)
	service := core.NewTaskService(repo, queuePublisher, cfg.QueueConfig, logger)
	handler := grpc.NewHandler(service)

	taskResultHandler := messaging.CreateTaskResultHandler(service, logger)
	taskResultConsumer := cm.NewMessageConsumer(
		queueManager,
		cfg.QueueConfig.TaskResultQueue,
		taskResultHandler,
		logger,
	)
	go func() {
		if err := taskResultConsumer.Start(); err != nil {
			logger.Error("Failed to start task result consumer", "error", err)
			return
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		logger.Error("Failed to listen", "error", err)
		return
	}

	grpcServer := grpc_server.NewServer(
		grpc_server.UnaryInterceptor(interceptors.LoggingInterceptor(logger)),
	)
	pb.RegisterTaskServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to start gRPC server", "error", err)
		return
	}
}
