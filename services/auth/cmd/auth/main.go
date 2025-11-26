package main

import (
	"log/slog"
	"os"

	"github.com/Anurag-S1ngh/carbon-backend/services/auth/internal/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := http.NewRouter()

	router.Run(":8000")
}
