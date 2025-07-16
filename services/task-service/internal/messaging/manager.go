package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQManager struct {
	connString  string
	connection  *amqp091.Connection
	notifyClose chan *amqp091.Error
}

func NewRabbitMQManager(connString string) *RabbitMQManager {
	return &RabbitMQManager{
		connString: connString,
	}
}

func (rm *RabbitMQManager) Connect() {
	var err error
	for {
		log.Println("Attempting to connect to RabbitMQ...")
		rm.connection, err = amqp091.Dial(rm.connString)
		if err == nil {
			log.Println("RabbitMQ connection successful")
			rm.notifyClose = make(chan *amqp091.Error)
			rm.connection.NotifyClose(rm.notifyClose)
			break
		}

		log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds. Error: %v", err)
		time.Sleep(5 * time.Second)
	}

	go rm.reconnect()
}

func (rm *RabbitMQManager) reconnect() {
	for {
		<-rm.notifyClose
		log.Println("RabbitMQ connection lost. Attempting to reconnect...")
		rm.Connect()
	}
}

func (rm *RabbitMQManager) GetChannel() (*amqp091.Channel, error) {
	if rm.connection == nil || rm.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}
	return rm.connection.Channel()
}

func (rm *RabbitMQManager) Close() {
	if rm.connection != nil && !rm.connection.IsClosed() {
		log.Println("Closing RabbitMQ connection")
		rm.connection.Close()
	}
}
