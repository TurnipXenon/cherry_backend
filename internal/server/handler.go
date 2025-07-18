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
	log.Println("Received webhook request from Todoist")

	// Get the client secret from environment variables
	clientSecret := os.Getenv("TODOIST_CLIENT_SECRET")
	if clientSecret == "" {
		log.Println("Warning: TODOIST_CLIENT_SECRET not set")
	}

	// Verify the request signature if client secret is available
	if clientSecret != "" {
		// Get the X-Todoist-Hmac-SHA256 header
		signature := r.Header.Get("X-Todoist-Hmac-SHA256")
		if signature == "" {
			log.Println("Error: Missing X-Todoist-Hmac-SHA256 header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Important: Restore the request body for later use
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		// Verify the signature
		if !verifyTodoistSignature(body, signature, clientSecret) {
			log.Println("Error: Invalid signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("Webhook signature verified successfully")
	} else {
		log.Println("Warning: Skipping signature verification as TODOIST_CLIENT_SECRET is not set")
	}

	// Parse the webhook payload
	var request apiv1.TodoistWebhookRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Process the webhook
	todoistService := &TodoistServiceImpl{}
	response, err := todoistService.ProcessWebhook(r.Context(), &request)
	if err != nil {
		log.Printf("Error processing webhook: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
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
	// Create request
	request := &apiv1.HealthCheckRequest{}

	// Process the health check
	healthService := &HealthServiceImpl{}
	response, err := healthService.Check(request)
	if err != nil {
		log.Printf("Error performing health check: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
