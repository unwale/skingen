package messaging

import (
	"fmt"
	"log"
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
}

func NewMessageConsumer(manager ChannelProvider, queueName string, handler MessageHandler) *MessageConsumer {
	consumerTag := generateConsumerTag()
	return &MessageConsumer{
		manager:     manager,
		consumerTag: consumerTag,
		queueName:   queueName,
		handler:     handler,
	}
}

func (c *MessageConsumer) Start() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
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
		return err
	}

	msgs, err := ch.Consume(
		c.queueName,
		c.consumerTag,
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
			err := c.handler(msg)
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

func (c *MessageConsumer) Shutdown() error {
	ch, err := c.manager.GetChannel()
	if err != nil {
		return err
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
	}()

	if err := ch.Cancel(c.consumerTag, false); err != nil {
		return fmt.Errorf("failed to cancel consumer: %w", err)
	}

	log.Printf("Consumer %s for queue %s has been shut down", c.consumerTag, c.queueName)
	return nil
}

func generateConsumerTag() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
		log.Printf("Failed to get hostname: %v", err)
	}

	pid := os.Getpid()

	tagUUID := uuid.NewString()

	tag := fmt.Sprintf("%s-pid%d-%s", hostname, pid, tagUUID[:8])

	return tag
}
