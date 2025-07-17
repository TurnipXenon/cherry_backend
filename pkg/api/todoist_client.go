package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TodoistClient is a client for the Todoist webhook API
type TodoistClient struct {
	generator *ClientGenerator
}

// NewTodoistClient creates a new Todoist client
func NewTodoistClient(baseURL string) *TodoistClient {
	return &TodoistClient{
		generator: NewClientGenerator(baseURL),
	}
}

// TodoistWebhookRequest represents a request to the Todoist webhook endpoint
type TodoistWebhookRequest struct {
	EventName string `json:"event_name"`
	UserID    string `json:"user_id"`
	EventData string `json:"event_data"`
	Version   string `json:"version"`
}

// TodoistWebhookResponse represents a response from the Todoist webhook endpoint
type TodoistWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ProcessWebhook sends a webhook notification to the Todoist webhook endpoint
func (c *TodoistClient) ProcessWebhook(req *TodoistWebhookRequest) (*TodoistWebhookResponse, error) {
	// Create request
	request := &Request{
		Method: http.MethodPost,
		Path:   "/webhooks/todoist",
		Body:   req,
	}

	// Execute request
	resp, err := c.generator.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to process webhook: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var response TodoistWebhookResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}
