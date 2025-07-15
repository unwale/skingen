package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port           string `env:"PORT"`
	TaskServiceUrl string `env:"TASK_SERVICE_URL"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
