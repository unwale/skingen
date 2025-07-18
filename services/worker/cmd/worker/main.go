package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	cm "github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/worker/internal/config"
	"github.com/unwale/skingen/services/worker/internal/core"
	"github.com/unwale/skingen/services/worker/internal/messaging"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration")
	}

	queueManager := cm.NewRabbitMQManager(cfg.RabbitMQUrl)
	queueManager.Connect()
	defer queueManager.Close()

	service := core.NewWorkerService()
	taskCommandHandler := messaging.CreateTaskCommandHandler(service)
	taskConsumer := cm.NewMessageConsumer(
		queueManager,
		cfg.QueueConfig.GenerationQueue,
		taskCommandHandler,
	)
	if err := taskConsumer.Start(); err != nil {
		log.Fatalf("Failed to start task consumer: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down worker service...")
	if err := taskConsumer.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown message consumer: %v", err)
	}
	queueManager.Close()
	log.Println("Worker service stopped gracefully")
}
