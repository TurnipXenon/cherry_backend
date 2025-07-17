package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server for the application
type Server struct {
	Router *mux.Router
}

// NewServer creates a new server instance
func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}

	// Register routes
	s.registerRoutes()

	return s
}

// registerRoutes sets up all the routes for the server
func (s *Server) registerRoutes() {
	// Register webhook handler
	s.Router.HandleFunc("/webhooks/todoist", s.TodoistWebhookHandler).Methods("POST")

	// Add a health check endpoint
	s.Router.HandleFunc("/health", s.HealthCheckHandler).Methods("GET")
}

// Run starts the HTTP server
func (s *Server) Run() error {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server starting on port %s...\n", port)
	return http.ListenAndServe(":"+port, s.Router)
}
