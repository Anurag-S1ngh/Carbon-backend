package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtSecret   string
	DatabaseURL string
	RabbitMQURL string
	Port        string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JwtSecret:   getEnv("JWT_SECRET", "StrongPassword12481jf"),
		Port:        getEnv("UPLOAD_SERVICE_PORT", "8003"),
	}
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
