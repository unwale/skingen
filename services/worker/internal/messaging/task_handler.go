package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/worker/internal/core"
)

func CreateTaskCommandHandler(service core.WorkerService) messaging.MessageHandler {
	return func(msg amqp091.Delivery) error {
		var command contracts.GenerateImageCommand
		if err := json.Unmarshal(msg.Body, &command); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		event, err := service.GenerateImage(context.Background(), &command)
		if err != nil {
			log.Printf("Failed to generate image: %v", err)
			return err
		}

		log.Printf("Generated image event: %+v", event)
		return nil
	}
}
