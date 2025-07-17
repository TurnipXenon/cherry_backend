package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"cherry_backend/internal/models"
)

func main() {
	// Define command line flags
	url := flag.String("url", "http://localhost:8080/webhooks/todoist", "URL to send the webhook to")
	secret := flag.String("secret", "", "Todoist client secret for signature verification")
	eventType := flag.String("event", "item:added", "Event type (item:added, item:updated, item:deleted, item:completed)")
	userID := flag.String("user", "12345", "User ID")
	flag.Parse()

	// Create a sample event data
	eventData := map[string]interface{}{
		"id":          "123456789",
		"content":     "Test task",
		"description": "This is a test task created by the webhook simulator",
		"due": map[string]interface{}{
			"date":        "2023-12-31",
			"is_recurring": false,
			"string":      "Dec 31",
		},
		"priority": 1,
	}

	// Create the webhook payload
	payload := models.TodoistWebhookPayload{
		EventName: *eventType,
		UserID:    *userID,
		EventData: eventData,
		Version:   "9",
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", *url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Calculate and set HMAC signature if secret is provided
	if *secret != "" {
		// Create a new HMAC by defining the hash type and the key
		h := hmac.New(sha256.New, []byte(*secret))

		// Write payload to the HMAC
		h.Write(jsonPayload)

		// Get the calculated signature
		signature := hex.EncodeToString(h.Sum(nil))

		// Set the signature header
		req.Header.Set("X-Todoist-Hmac-SHA256", signature)
		fmt.Println("Added signature header with value:", signature)
	} else {
		fmt.Println("Warning: No secret provided, skipping signature calculation")
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)

	// Print success message
	fmt.Printf("Successfully sent %s webhook to %s\n", *eventType, *url)
}
