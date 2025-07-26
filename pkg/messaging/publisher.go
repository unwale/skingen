package messaging

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	manager ChannelProvider
	logger  *slog.Logger
}

func NewRabbitMQPublisher(manager ChannelProvider, logger *slog.Logger) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		manager: manager,
		logger:  logger,
	}
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, body []byte, queueName, correlationID string) error {
	ch, err := p.manager.GetChannel()
	if err != nil {
		return err
	}
	defer func() {
		if err := ch.Close(); err != nil {
			p.logger.Error("Failed to close channel", "error", err)
		}
	}()

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          body,
		},
	)
	if err != nil {
		return err
	}

	p.logger.Info("Message published successfully",
		"queue", queueName,
		"correlation_id", correlationID,
	)
	return nil
}
