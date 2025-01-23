package config

import (
	"effective-mobile-task/internal/database/postgres"
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	ServerPort string `env:"SERVER_PORT"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)

	if cfg == (Config{}) {
		return nil, errors.New("config is empty")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
