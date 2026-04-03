package config

type RedisConfig struct {
	// URL string `env:"REDIS_URL,required"`
	Address string `env:"REDIS_ADDR,required"`
}
