package messaging

import (
	"context"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	manager ChannelProvider
}

func NewRabbitMQPublisher(manager ChannelProvider) *RabbitMQPublisher {
	return &RabbitMQPublisher{manager: manager}
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, body []byte, queueName, correlationID string) error {
	ch, err := p.manager.GetChannel()
	if err != nil {
		return err
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
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

	log.Printf("Successfully published message to queue: %s", queueName)
	return nil
}
