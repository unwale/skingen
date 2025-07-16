package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/services/task-service/internal/config"
	"github.com/unwale/skingen/services/task-service/internal/core"
)

type MessageConsumer struct {
	manager     ChannelProvider
	service     core.TaskService
	queueConfig config.QueueConfig
}

func NewMessageConsumer(manager ChannelProvider, service core.TaskService, queueCfg config.QueueConfig) *MessageConsumer {
	return &MessageConsumer{
		manager:     manager,
		service:     service,
		queueConfig: queueCfg,
	}
}

func (c *MessageConsumer) Start() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queueConfig.GenerateImageQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := c.handleMessage(msg)
			if err != nil {
				log.Printf("Error handling message: %v", err)
				if nackErr := ch.Nack(msg.DeliveryTag, false, true); nackErr != nil {
					log.Printf("Error nack'ing message: %v", nackErr)
				}
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}

func (c *MessageConsumer) handleMessage(msg amqp091.Delivery) error {
	var event contracts.GenerateImageEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	_, err := c.service.ProcessTaskResult(
		context.Background(),
		event,
	)
	return err
}
