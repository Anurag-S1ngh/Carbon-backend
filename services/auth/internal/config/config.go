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
}

func Load() *Config {
	godotenv.Load()

	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres")
	resendAPIKey := getEnv("RESEND_API_KEY", "")
	redisURL := getEnv("REDIS_URL", "")
	otpExpirySeconds := getEnv("OTP_EXPIRY_SECONDS", "600")
	refreshTokenExpiryHours := getEnv("REFRESH_TOKEN_EXPIRY_HOURS", "48")
	accessTokenExpiryHours := getEnv("ACCESS_TOKEN_EXPIRY_HOURS", "4")
	jwtSecret := getEnv("JWT_SECRET", "strongPassword")

	return &Config{
		DatabaseURL:             databaseURL,
		ResendAPIKey:            resendAPIKey,
		RedisURL:                redisURL,
		OTPExpirySeconds:        otpExpirySeconds,
		RefreshTokenExpiryHours: refreshTokenExpiryHours,
		AccessTokenExpiryHours:  accessTokenExpiryHours,
		JwtSecret:               jwtSecret,
	}
}

func getEnv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return v
}
