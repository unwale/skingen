package messaging

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(msg amqp091.Delivery) error

type MessageConsumer struct {
	manager   ChannelProvider
	queueName string
	handler   MessageHandler
}

func NewMessageConsumer(manager ChannelProvider, queueName string, handler MessageHandler) *MessageConsumer {
	return &MessageConsumer{
		manager:   manager,
		queueName: queueName,
		handler:   handler,
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
