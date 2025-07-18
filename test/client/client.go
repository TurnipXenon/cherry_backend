package main

import (
	"fmt"
	"log"

	"cherry_backend/pkg/api"
)

func main() {
	// Test the Todoist client
	testTodoistClient()

	// Test the Health client
	testHealthClient()
}

func testTodoistClient() {
	fmt.Println("Testing Todoist client...")

	// Create a Todoist client
	todoistClient := api.NewTodoistClient("http://localhost:8080")

	// Set the client secret for HMAC signature calculation
	// This should match the TODOIST_CLIENT_SECRET environment variable on the server
	todoistClient.SetSecret("test_secret")

	// Create a webhook request
	request := &api.TodoistWebhookRequest{
		EventName: "item:added",
		UserID:    "123456",
		EventData: `{"item_id": "789"}`,
		Version:   "1.0",
	}

	// Process the webhook
	response, err := todoistClient.ProcessWebhook(request)
	if err != nil {
		log.Fatalf("Failed to process webhook: %v", err)
	}

	// Print the response
	fmt.Printf("Success: %v, Message: %s\n", response.Success, response.Message)
}

func testHealthClient() {
	fmt.Println("Testing Health client...")

	// Create a health check client
	healthClient := api.NewHealthClient("http://localhost:8080")

	// Perform a health check
	response, err := healthClient.Check()
	if err != nil {
		log.Fatalf("Failed to perform health check: %v", err)
	}

	// Print the response
	fmt.Printf("Status: %s\n", response.Status)
}
