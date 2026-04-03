package config

type ServerConfig struct {
	JwtSecret string `env:"JWT_SECRET,required"`
	Port      int    `env:"PORT,required"`
}
