package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv    string `env:"APP_ENV" envDefault:"development"`
	Port      int    `env:"PORT" envDefault:"9090"`
	RedisAddr string `env:"REDIS_ADDR,required"`
	Debug     bool   `env:"DEBUG" envDefault:"true"`
}

func Load() (*Config, error) {
	loadDotEnv()

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}

func validate(cfg *Config) error {
	// TODO: do validation
	return nil
}
