package main

import (
	"context"
	"log/slog"
	"os"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/connection"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/email"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/redis"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/token/jwt"
	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/config"
	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/http"
	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/service"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	cfg := config.Load()

	goth.UseProviders(github.New(cfg.GithubClientID, cfg.GithubClientSecret, cfg.GithubCallbackURL, "user:email"))

	ctx := context.Background()

	dbConn, err := db.NewDatabaseConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("error while connecting to database", "error", err)
		return
	}
	err = dbConn.Ping(ctx)
	if err != nil {
		logger.Error("error while ping database", "error", err)
		return
	}
	defer dbConn.Close(ctx)
	dbQueries := db.DatabaseQueries(dbConn)

	emailConfig := email.NewEmailConfig(cfg.ResendAPIKey, logger)
	jwtConfig := jwt.NewJwtConfig(cfg.JwtSecret, logger)
	redisConfig, err := redis.NewRedisConfig(cfg.RedisURL, logger)
	if err != nil {
		logger.Error("error while creatinge redis config", "error", err)
		return
	}

	authService := service.NewAuthService(
		dbQueries,
		cfg.OTPExpirySeconds,
		cfg.RefreshTokenExpiryHours,
		cfg.AccessTokenExpiryHours,
		emailConfig,
		redisConfig,
		jwtConfig,
		logger)

	router := http.NewRouter(authService, cfg.RedirectURL)

	logger.Info("auth service started on 8000", "info", 8000)
	err = router.Run(":8000")
	if err != nil {
		logger.Error("error while starting the server", "error", err)
	}
}
