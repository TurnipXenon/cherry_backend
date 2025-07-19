package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	apiv1 "cherry_backend/pkg/api/v1"
)

// TodoistWebhookHandler processes incoming webhook notifications from Todoist
func (s *Server) TodoistWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new TodoistServiceImpl with a logger
	todoistService, err := NewTodoistServiceImpl()
	if err != nil {
		log.Printf("Error creating TodoistServiceImpl: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use the logger for all logging
	logger := todoistService.Logger
	logger.Info("Received webhook request from Todoist")

	// Get the client secret from environment variables
	clientSecret := os.Getenv("TODOIST_CLIENT_SECRET")
	if clientSecret == "" {
		logger.Warn("TODOIST_CLIENT_SECRET not set")
	}

	// Verify the request signature if client secret is available
	if clientSecret != "" {
		// Get the X-Todoist-Hmac-SHA256 header
		signature := r.Header.Get("X-Todoist-Hmac-SHA256")
		if signature == "" {
			logger.Error("Missing X-Todoist-Hmac-SHA256 header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading request body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Important: Restore the request body for later use
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// Verify the signature
		if !verifyTodoistSignature(body, signature, clientSecret) {
			logger.Error("Invalid signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		logger.Info("Webhook signature verified successfully")
	} else {
		logger.Warn("Skipping signature verification as TODOIST_CLIENT_SECRET is not set")
	}

	// Parse the webhook payload
	var request apiv1.TodoistWebhookRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Error("Error parsing webhook payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Process the webhook
	response, err := todoistService.ProcessWebhook(r.Context(), &request)
	if err != nil {
		logger.Error("Error processing webhook: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// verifyTodoistSignature verifies the HMAC-SHA256 signature from Todoist
func verifyTodoistSignature(payload []byte, signature, secret string) bool {
	// Create a new HMAC by defining the hash type and the key
	h := hmac.New(sha256.New, []byte(secret))

	// Write payload to the HMAC
	h.Write(payload)

	// Get the calculated signature
	calculatedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare the calculated signature with the provided one
	return hmac.Equal([]byte(calculatedSignature), []byte(signature))
}

// HealthCheckHandler performs a health check
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new TodoistServiceImpl with a logger (we'll use this logger for health checks too)
	todoistService, err := NewTodoistServiceImpl()
	if err != nil {
		log.Printf("Error creating TodoistServiceImpl for health check: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use the logger for all logging
	logger := todoistService.Logger
	logger.Info("Received health check request")

	// Create request
	request := &apiv1.HealthCheckRequest{}

	// Process the health check
	healthService := &HealthServiceImpl{}
	response, err := healthService.Check(request)
	if err != nil {
		logger.Error("Error performing health check: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Health check successful")

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
