package server

import (
	"context"
	"log"

	"cherry_backend/internal/logging"
	apiv1 "cherry_backend/pkg/api/v1"
)

// TodoistServiceImpl implements the TodoistService interface
type TodoistServiceImpl struct {
	apiv1.UnimplementedTodoistServiceServer
	Logger logging.Logger
}

// NewTodoistServiceImpl creates a new TodoistServiceImpl with a logger
func NewTodoistServiceImpl() (*TodoistServiceImpl, error) {
	logger, err := logging.NewLogger()
	if err != nil {
		return nil, err
	}

	return &TodoistServiceImpl{
		Logger: logger,
	}, nil
}

// ProcessWebhook processes incoming webhook notifications from Todoist
func (s *TodoistServiceImpl) ProcessWebhook(ctx context.Context, request *apiv1.TodoistWebhookRequest) (*apiv1.TodoistWebhookResponse, error) {
	// Use the standard log package as a fallback if logger is not initialized
	if s.Logger == nil {
		log.Printf("Warning: Logger not initialized, using standard log package")
		log.Printf("Processing webhook event: %s", request.EventName)
	} else {
		s.Logger.Info("Processing webhook event: %s", request.EventName)
	}

	// Here you would add your business logic to handle different event types
	// For example:
	switch request.EventName {
	case "item:added":
		if s.Logger != nil {
			s.Logger.Info("Item added by user %s", request.UserId)
		} else {
			log.Printf("Item added by user %s", request.UserId)
		}
		// Handle item added event
	case "item:updated":
		if s.Logger != nil {
			s.Logger.Info("Item updated by user %s", request.UserId)
		} else {
			log.Printf("Item updated by user %s", request.UserId)
		}
		// Handle item updated event
	case "item:deleted":
		if s.Logger != nil {
			s.Logger.Info("Item deleted by user %s", request.UserId)
		} else {
			log.Printf("Item deleted by user %s", request.UserId)
		}
		// Handle item deleted event
	case "item:completed":
		if s.Logger != nil {
			s.Logger.Info("Item completed by user %s", request.UserId)
		} else {
			log.Printf("Item completed by user %s", request.UserId)
		}
		// Handle item completed event
	default:
		if s.Logger != nil {
			s.Logger.Warn("Unhandled event type: %s", request.EventName)
		} else {
			log.Printf("Unhandled event type: %s", request.EventName)
		}
	}

	// Return success response
	return &apiv1.TodoistWebhookResponse{
		Success: true,
		Message: "Webhook received",
	}, nil
}
