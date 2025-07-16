package messaging

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type AMQPChannel interface {
	Close() error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error)
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error
}

type ChannelProvider interface {
	GetChannel() (AMQPChannel, error)
}
