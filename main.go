package main

import (
	"log"

	"github.com/joho/godotenv"

	"cherry_backend/internal/server"
)

func main() {
	// Load environment variables from .env file if it exists
	err := godotenv.Load("configs/local.env")
	if err != nil {
		log.Println("Warning: configs/local.env file not found. Using environment variables.")
	}

	// Create a new server
	s := server.NewServer()

	// Start the server
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
