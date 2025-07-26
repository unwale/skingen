package messaging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(msg amqp091.Delivery) error

type MessageConsumer struct {
	manager     ChannelProvider
	consumerTag string
	queueName   string
	handler     MessageHandler
	logger      *slog.Logger
}

func NewMessageConsumer(manager ChannelProvider, queueName string, handler MessageHandler, logger *slog.Logger) *MessageConsumer {
	consumerTag := generateConsumerTag()
	logger = logger.With(slog.String("consumer_tag", consumerTag), slog.String("queue_name", queueName))
	logger.Info("Creating new message consumer", "consumer_tag", consumerTag, "queue_name", queueName)
	return &MessageConsumer{
		manager:     manager,
		consumerTag: consumerTag,
		queueName:   queueName,
		handler:     handler,
		logger:      logger,
	}
}

func (c *MessageConsumer) Start() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
		c.logger.Error("Failed to get channel", "error", err)
		return err
	}

	_, err = ch.QueueDeclare(
		c.queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.logger.Error("Failed to declare queue", "error", err, "queue_name", c.queueName)
		return err
	}

	msgs, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to start consuming messages", "error", err, "queue_name", c.queueName)
		return err
	}

	go func() {
		for msg := range msgs {
			logger := c.logger.With("correlation_id", msg.CorrelationId)
			err := c.handler(msg)
			if err != nil {
				logger.Error("Error processing message", "error", err)
				if nackErr := ch.Nack(msg.DeliveryTag, false, true); nackErr != nil {
					logger.Error("Failed to nack message", "error", nackErr)
				}
				continue
			}
			if err := msg.Ack(false); err != nil {
				logger.Error("Failed to ack message", "error", err)
			} else {
				logger.Info("Message processed successfully")
			}
		}
	}()

	return nil
}

func (c *MessageConsumer) Shutdown() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
		c.logger.Error("Failed to get channel for shutdown", "error", err)
		return err
	}
	defer func() {
		if err := ch.Close(); err != nil {
			c.logger.Error("Failed to close channel during shutdown", "error", err)
		}
	}()

	if err := ch.Cancel(c.consumerTag, false); err != nil {
		c.logger.Error("Failed to cancel consumer", "error", err)
		return err
	}

	c.logger.Info("Message consumer shutdown successfully", "consumer_tag", c.consumerTag, "queue_name", c.queueName)
	return nil
}

func generateConsumerTag() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
		slog.Warn("Failed to get hostname, using 'unknown-host'", "error", err)
	}

	pid := os.Getpid()

	tagUUID := uuid.NewString()

	tag := fmt.Sprintf("%s-pid%d-%s", hostname, pid, tagUUID[:8])

	return tag
}
