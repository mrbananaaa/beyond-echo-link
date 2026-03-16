package config

import (
	"os"

	"github.com/joho/godotenv"
)

func loadDotEnv() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	if env == "production" {
		return
	}

	_ = godotenv.Load()

	_ = godotenv.Load(".env.local")
}
