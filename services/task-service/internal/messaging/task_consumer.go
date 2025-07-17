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

type TaskResultConsumer struct {
	manager     ChannelProvider
	service     core.TaskService
	queueConfig config.QueueConfig
}

func NewTaskResultConsumer(manager ChannelProvider, service core.TaskService, queueCfg config.QueueConfig) *TaskResultConsumer {
	return &TaskResultConsumer{
		manager:     manager,
		service:     service,
		queueConfig: queueCfg,
	}
}

func (c *TaskResultConsumer) Start() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		c.queueConfig.TaskResultQueue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queueConfig.TaskResultQueue,
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
			if err := msg.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}
		}
	}()

	return nil
}

func (c *TaskResultConsumer) handleMessage(msg amqp091.Delivery) error {
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
