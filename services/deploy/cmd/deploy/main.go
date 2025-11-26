package main

import "github.com/Anurag-S1ngh/carbon-backend/services/deploy/internal/http"

func main() {
	router := http.NewRouter()

	router.Run(":8001")
}
