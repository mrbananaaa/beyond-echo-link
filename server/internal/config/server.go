package config

type ServerConfig struct {
	Port int `env:"PORT,required"`
}
