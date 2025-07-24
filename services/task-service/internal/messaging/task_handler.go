package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/task-service/internal/core"
)

func CreateTaskResultHandler(service core.TaskService) messaging.MessageHandler {
	return func(msg amqp091.Delivery) error {
		var event contracts.GenerateImageEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		task, err := service.ProcessTaskResult(context.Background(), event)
		if err != nil {
			log.Printf("Failed to process task result: %v", err)
			return err
		}

		log.Printf("Processed task result: %+v", task)
		return nil
	}
}
