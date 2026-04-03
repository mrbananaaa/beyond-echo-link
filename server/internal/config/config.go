package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`

	Server ServerConfig
	DB     DatabaseConfig
	Log    LogConfig
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
	if cfg.Server.Port <= 0 {
		return fmt.Errorf("invalid server port")
	}

	if cfg.DB.URL == "" {
		return fmt.Errorf("database url is required")
	}

	return nil
}
