package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port             string `env:"PORT,required"`
	ServiceName      string `env:"SERVICE_NAME" envDefault:"task-service"`
	LoggingLevel     string `env:"LOGGING_LEVEL" envDefault:"info"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     string `env:"POSTGRES_PORT,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	RabbitMQUrl      string `env:"RABBITMQ_URL,required"`
	QueueConfig      QueueConfig
}

type QueueConfig struct {
	GenerateImageQueue string `env:"GENERATE_IMAGE_QUEUE,required,notEmpty"`
	TaskResultQueue    string `env:"TASK_RESULT_QUEUE,required,notEmpty"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
