package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Port string `env:"PORT"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
