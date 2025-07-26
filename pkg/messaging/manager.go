package messaging

import (
	"errors"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQManager struct {
	connString  string
	connection  *amqp091.Connection
	notifyClose chan *amqp091.Error
	logger      *slog.Logger
}

func NewRabbitMQManager(connString string, logger *slog.Logger) *RabbitMQManager {
	return &RabbitMQManager{
		connString: connString,
		logger:     logger,
	}
}

func (rm *RabbitMQManager) Connect() {
	var err error
	for {
		rm.logger.Info("Attempting to connect to RabbitMQ", "connection_string", rm.connString)
		rm.connection, err = amqp091.Dial(rm.connString)
		if err == nil {
			rm.logger.Info("Connected to RabbitMQ successfully")
			rm.notifyClose = make(chan *amqp091.Error)
			rm.connection.NotifyClose(rm.notifyClose)
			break
		}

		rm.logger.Warn("Failed to connect to RabbitMQ", "error", err)
		time.Sleep(5 * time.Second)
	}

	go rm.reconnect()
}

func (rm *RabbitMQManager) reconnect() {
	for {
		<-rm.notifyClose
		rm.logger.Warn("RabbitMQ connection closed, attempting to reconnect")
		rm.Connect()
	}
}

func (rm *RabbitMQManager) GetChannel() (AMQPChannel, error) {
	if rm.connection == nil || rm.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}
	return rm.connection.Channel()
}

func (rm *RabbitMQManager) Close() {
	if rm.connection != nil && !rm.connection.IsClosed() {
		rm.logger.Info("Closing RabbitMQ connection")
		rm.connection.Close() //nolint:errcheck
	}
}
