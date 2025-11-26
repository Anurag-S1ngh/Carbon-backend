package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
}

func NewConfig() *Config {
	godotenv.Load()

	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres")

	return &Config{
		DBURL: databaseURL,
	}
}

func getEnv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return v
}
