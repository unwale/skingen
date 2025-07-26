package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port           string `env:"PORT,required"`
	ServiceName    string `env:"SERVICE_NAME" envDefault:"gateway"`
	LoggingLevel   string `env:"LOGGING_LEVEL" envDefault:"info"`
	TaskServiceUrl string `env:"TASK_SERVICE_URL,required"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
