package main

import "github.com/Anurag-S1ngh/carbon-backend/services/upload/internal/http"

func main() {
	router := http.NewRouter()

	router.Run(":8003")
}
