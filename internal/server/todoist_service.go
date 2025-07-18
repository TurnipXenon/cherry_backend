package server

import (
	"context"
	"log"

	apiv1 "cherry_backend/pkg/api/v1"
)

// TodoistServiceImpl implements the TodoistService interface
type TodoistServiceImpl struct {
	apiv1.UnimplementedTodoistServiceServer
	// Add any dependencies here
}

// ProcessWebhook processes incoming webhook notifications from Todoist
func (s *TodoistServiceImpl) ProcessWebhook(ctx context.Context, request *apiv1.TodoistWebhookRequest) (*apiv1.TodoistWebhookResponse, error) {
	log.Printf("Processing webhook event: %s", request.EventName)

	// Here you would add your business logic to handle different event types
	// For example:
	switch request.EventName {
	case "item:added":
		log.Printf("Item added by user %s", request.UserId)
		// Handle item added event
	case "item:updated":
		log.Printf("Item updated by user %s", request.UserId)
		// Handle item updated event
	case "item:deleted":
		log.Printf("Item deleted by user %s", request.UserId)
		// Handle item deleted event
	case "item:completed":
		log.Printf("Item completed by user %s", request.UserId)
		// Handle item completed event
	default:
		log.Printf("Unhandled event type: %s", request.EventName)
	}

	// Return success response
	return &apiv1.TodoistWebhookResponse{
		Success: true,
		Message: "Webhook received",
	}, nil
}
