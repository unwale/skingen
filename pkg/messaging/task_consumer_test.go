package messaging

import (
	"context"
	"testing"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChannelProvider struct {
	mock.Mock
}

func (m *mockChannelProvider) GetChannel() (AMQPChannel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(AMQPChannel), args.Error(1)
}

type mockAMQPChannel struct {
	mock.Mock
}

func (m *mockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	callArgs := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return callArgs.Get(0).(amqp091.Queue), callArgs.Error(1)
}

func (m *mockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	callArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return callArgs.Get(0).(<-chan amqp091.Delivery), callArgs.Error(1)
}

func (m *mockAMQPChannel) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	callArgs := m.Called(ctx, exchange, key, mandatory, immediate, msg)
	return callArgs.Error(0)
}

func (m *mockAMQPChannel) Ack(tag uint64, multiple bool) error {
	args := m.Called(tag, multiple)
	return args.Error(0)
}

func (m *mockAMQPChannel) Nack(tag uint64, multiple, requeue bool) error {
	args := m.Called(tag, multiple, requeue)
	return args.Error(0)
}

func (m *mockAMQPChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestStartConsuming(t *testing.T) {
	manager := &mockChannelProvider{}
	queueName := "queue"
	handler := func(msg amqp091.Delivery) error {
		return nil
	}

	consumer := NewTaskResultConsumer(manager, queueName, handler)

	mockChannel := new(mockAMQPChannel)
	mockQueue := make(<-chan amqp091.Delivery)
	manager.On("GetChannel").Return(mockChannel, nil)
	mockChannel.On("QueueDeclare", queueName, true, false, false, false, mock.Anything).Return(amqp091.Queue{}, nil)
	mockChannel.On("Consume", queueName, "", true, false, false, false, mock.Anything).Return(mockQueue, nil)

	err := consumer.Start()

	assert.NoError(t, err)
	mockChannel.AssertExpectations(t)
}
