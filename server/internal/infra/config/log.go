package config

type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"debug"`
}
