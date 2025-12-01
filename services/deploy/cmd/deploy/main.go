package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Anurag-S1ngh/carbon-backend/services/deploy/internal/config"
	"github.com/Anurag-S1ngh/carbon-backend/services/deploy/internal/http"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	router := http.NewRouter()
	port := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("upload service started", "port", port)
	err := router.Run(port)
	if err != nil {
		logger.Error("error while starting the server", "error", err)
	}
}
