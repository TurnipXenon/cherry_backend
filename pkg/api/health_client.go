package api

import (
	"fmt"
	"net/http"
)

// HealthClient is a client for the health check API
type HealthClient struct {
	generator *ClientGenerator
}

// NewHealthClient creates a new health check client
func NewHealthClient(baseURL string) *HealthClient {
	return &HealthClient{
		generator: NewClientGenerator(baseURL),
	}
}

// HealthCheckRequest represents a request to the health check endpoint
type HealthCheckRequest struct {
	// No parameters needed for a simple health check
}

// HealthCheckResponse represents a response from the health check endpoint
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// Check performs a health check
func (c *HealthClient) Check() (*HealthCheckResponse, error) {
	// Create request
	request := &Request{
		Method: http.MethodGet,
		Path:   "/health",
	}

	// Execute request
	resp, err := c.generator.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to perform health check: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create response
	response := &HealthCheckResponse{
		Status: string(resp.Body),
	}

	return response, nil
}
