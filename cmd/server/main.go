package main

import (
	"log"

	"github.com/org/ranas-bdi-backend/internal/adapters/http"
)

func main() {
	handler := http.NewHandler()
	router := handler.GetRouter()

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
