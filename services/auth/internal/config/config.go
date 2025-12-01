package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL             string
	RedisURL                string
	ResendAPIKey            string
	OTPExpirySeconds        string
	RefreshTokenExpiryHours string
	AccessTokenExpiryHours  string
	JwtSecret               string
	GithubClientID          string
	GithubClientSecret      string
	GithubCallbackURL       string
	RedirectURL             string
	Port                    string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:             getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres"),
		ResendAPIKey:            getEnv("RESEND_API_KEY", ""),
		RedisURL:                getEnv("REDIS_URL", ""),
		OTPExpirySeconds:        getEnv("OTP_EXPIRY_SECONDS", "600"),
		RefreshTokenExpiryHours: getEnv("REFRESH_TOKEN_EXPIRY_HOURS", "48"),
		AccessTokenExpiryHours:  getEnv("ACCESS_TOKEN_EXPIRY_HOURS", "4"),
		JwtSecret:               getEnv("JWT_SECRET", "strongPassword"),
		GithubClientID:          getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:      getEnv("GITHUB_CLIENT_SECRET", ""),
		GithubCallbackURL:       getEnv("GIHTUB_CALLBACK_URL", ""),
		RedirectURL:             getEnv("REDIRECT_URL", ""),
		Port:                    getEnv("AUTH_SERVICE_PORT", "8000"),
	}
}

func getEnv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return v
}
