package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port           string `env:"PORT,required"`
	RabbitMQUrl    string `env:"RABBITMQ_URL,required"`
	ModelServerUrl string `env:"MODEL_SERVER_URL,required"`
	QueueConfig    QueueConfig
	S3Config       S3Config
}

type QueueConfig struct {
	GenerationQueue string `env:"GENERATE_IMAGE_QUEUE,required"`
	TaskResultQueue string `env:"TASK_RESULT_QUEUE,required"`
}

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT,required"`
	AccessKey string `env:"S3_ACCESS_KEY,required"`
	SecretKey string `env:"S3_SECRET_KEY,required"`
	Bucket    string `env:"S3_BUCKET,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
