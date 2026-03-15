package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string `env:"APP_ENV" envDefault:"development"`
	Port      int    `env:"PORT" envDefault:"9090"`
	RedisAddr string `env:"REDIS_ADDR,required"`
	Debug     bool   `env:"DEBUG" envDefault:"true"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
