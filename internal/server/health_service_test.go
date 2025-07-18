package server

import (
	"testing"

	apiv1 "cherry_backend/pkg/api/v1"
)

func TestHealthCheck(t *testing.T) {
	// Create a new HealthServiceImpl
	service := &HealthServiceImpl{}

	// Create a request
	request := &apiv1.HealthCheckRequest{}

	// Call the Check method
	response, err := service.Check(request)

	// Check for errors
	if err != nil {
		t.Fatalf("Check returned an error: %v", err)
	}

	// Check the response
	if response == nil {
		t.Fatal("Check returned a nil response")
	}

	// Check the status
	expectedStatus := "OK"
	if response.Status != expectedStatus {
		t.Errorf("Check() status = %s, want %s", response.Status, expectedStatus)
	}
}