package config

import (
	"effective-mobile-task/internal/database/postgres"
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	ServerPort string `env:"SERVER_PORT"`
	ServerHost string `env:"SERVER_HOST"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}

	if cfg == (Config{}) {
		return nil, errors.New("config is empty")
	}

	return &cfg, nil
}
