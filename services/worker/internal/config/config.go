package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port        string `env:"PORT,required"`
	RabbitMQUrl string `env:"RABBITMQ_URL,required"`
	QueueConfig QueueConfig
}

type QueueConfig struct {
	GenerationQueue string `env:"GENERATE_IMAGE_QUEUE,required"`
	TaskResultQueue string `env:"TASK_RESULT_QUEUE,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
