package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtSecret   string
	DatabaseURL string
	RabbitMQURL string
}

func Load() *Config {
	godotenv.Load()

	jwtSecret := getEnv("JWT_SECRET", "StrongPassword12481jf")
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres")
	rabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	return &Config{
		JwtSecret:   jwtSecret,
		DatabaseURL: databaseURL,
		RabbitMQURL: rabbitMQURL,
	}
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
