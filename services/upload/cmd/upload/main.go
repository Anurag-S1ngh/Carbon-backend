package main

import (
	"context"
	"log/slog"
	"os"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/connection"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/middleware"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/rabbitmq"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/token/jwt"
	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/config"
	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/http"
	"github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/service"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	cfg := config.Load()

	dbConn, err := db.NewDatabaseConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("error while ping database", "error", err)
		return
	}
	err = dbConn.Ping(ctx)
	if err != nil {
		logger.Error("error while ping database", "error", err)
		return
	}
	defer dbConn.Close(ctx)
	dbQueries := db.DatabaseQueries(dbConn)

	rabbitmqConfig, err := rabbitmq.Connect(cfg.RabbitMQURL)
	if err != nil {
		logger.Error("error while connecting rabbitmq", "error", err)
		return
	}

	jwtConfig := jwt.NewJwtConfig(cfg.JwtSecret, logger)
	authMiddlewareConfig := middleware.NewAuthMiddlewareConfig(jwtConfig)
	gitService := service.NewGitService(dbQueries, logger)
	uploadServcie := service.NewUploadService(dbQueries, rabbitmqConfig, logger)

	router := http.NewRouter(gitService, uploadServcie, authMiddlewareConfig)

	err = router.Run(":8003")
	if err != nil {
		logger.Error("error while runnin the server", "error", err)
	}
}
