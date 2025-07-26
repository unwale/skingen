package messaging

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChannelProvider struct {
	mock.Mock
}

func (m *MockChannelProvider) GetChannel() (AMQPChannel, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(AMQPChannel), args.Error(1)
}

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAMQPChannel) Cancel(consumer string, noWait bool) error {
	args := m.Called(consumer, noWait)
	return args.Error(0)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	callArgs := m.Called(name, durable, autoDelete, exclusive, noWait, args)
	return callArgs.Get(0).(amqp091.Queue), callArgs.Error(1)
}

func (m *MockAMQPChannel) PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	callArgs := m.Called(ctx, exchange, key, mandatory, immediate, msg)
	return callArgs.Error(0)
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	callArgs := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	return callArgs.Get(0).(<-chan amqp091.Delivery), callArgs.Error(1)
}

func (m *MockAMQPChannel) Ack(tag uint64, multiple bool) error {
	args := m.Called(tag, multiple)
	return args.Error(0)
}

func (m *MockAMQPChannel) Nack(tag uint64, multiple, requeue bool) error {
	args := m.Called(tag, multiple, requeue)
	return args.Error(0)
}

func TestPublish(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockProvider := new(MockChannelProvider)
		mockChannel := new(MockAMQPChannel)

		publisher := NewRabbitMQPublisher(mockProvider, logger)

		ctx := context.Background()
		queueName := "test_queue"
		correlationID := uuid.New().String()
		body := []byte(`{"message": "hello"}`)

		mockProvider.On("GetChannel").Return(mockChannel, nil)
		mockChannel.On("QueueDeclare", queueName, true, false, false, false, mock.Anything).Return(amqp091.Queue{}, nil)
		mockChannel.On("PublishWithContext", ctx, "", queueName, false, false, amqp091.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          body,
		}).Return(nil)
		mockChannel.On("Close").Return(nil)

		err := publisher.Publish(ctx, body, queueName, correlationID)

		assert.NoError(t, err)
		mockProvider.AssertExpectations(t)
		mockChannel.AssertExpectations(t)
	})

	t.Run("failure on GetChannel", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockProvider := new(MockChannelProvider)
		publisher := NewRabbitMQPublisher(mockProvider, logger)

		expectedErr := errors.New("could not get channel")
		mockProvider.On("GetChannel").Return(nil, expectedErr)

		err := publisher.Publish(context.Background(), []byte("test"), "test_queue", "123")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockProvider.AssertExpectations(t)
	})

	t.Run("failure on QueueDeclare", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockProvider := new(MockChannelProvider)
		mockChannel := new(MockAMQPChannel)
		publisher := NewRabbitMQPublisher(mockProvider, logger)

		expectedErr := errors.New("permission denied for queue")
		mockProvider.On("GetChannel").Return(mockChannel, nil)
		mockChannel.On("QueueDeclare", "test_queue", true, false, false, false, mock.Anything).Return(amqp091.Queue{}, expectedErr)
		mockChannel.On("Close").Return(nil)

		err := publisher.Publish(context.Background(), []byte("test"), "test_queue", "123")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockProvider.AssertExpectations(t)
		mockChannel.AssertExpectations(t)
	})

	t.Run("failure on PublishWithContext", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))
		mockProvider := new(MockChannelProvider)
		mockChannel := new(MockAMQPChannel)
		publisher := NewRabbitMQPublisher(mockProvider, logger)

		expectedErr := errors.New("publish failed")
		mockProvider.On("GetChannel").Return(mockChannel, nil)
		mockChannel.On("QueueDeclare", "test_queue", true, false, false, false, mock.Anything).Return(amqp091.Queue{}, nil)
		mockChannel.On("PublishWithContext", mock.Anything, "", "test_queue", false, false, mock.Anything).Return(expectedErr)
		mockChannel.On("Close").Return(nil)

		err := publisher.Publish(context.Background(), []byte("test"), "test_queue", "123")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockProvider.AssertExpectations(t)
		mockChannel.AssertExpectations(t)
	})
}
