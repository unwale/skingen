package main

import (
	"fmt"
	"log"
	"net"

	grpc_server "google.golang.org/grpc"

	cm "github.com/unwale/skingen/pkg/messaging"
	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"github.com/unwale/skingen/services/task-service/internal/api/grpc"
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

	queueManager := cm.NewRabbitMQManager(cfg.RabbitMQUrl)
	queueManager.Connect()
	defer queueManager.Close()

	db, err := database.NewConnection(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queuePublisher := cm.NewRabbitMQPublisher(queueManager)
	repo := repository.NewTaskRepository(db)
	service := core.NewTaskService(repo, queuePublisher, cfg.QueueConfig)
	handler := grpc.NewHandler(service)

	taskResultHandler := messaging.CreateTaskResultHandler(service)
	taskResultConsumer := cm.NewMessageConsumer(
		queueManager,
		cfg.QueueConfig.GenerateImageQueue,
		taskResultHandler,
	)
	go func() {
		if err := taskResultConsumer.Start(); err != nil {
			log.Fatalf("Failed to start message consumer: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc_server.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, handler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
