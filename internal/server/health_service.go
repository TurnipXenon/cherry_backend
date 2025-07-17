package server

import (
	"cherry_backend/pkg/api/v1"
)

// HealthServiceImpl implements the HealthService interface
type HealthServiceImpl struct {
	// Add any dependencies here
}

// Check performs a health check
func (s *HealthServiceImpl) Check(request *apiv1.HealthCheckRequest) (*apiv1.HealthCheckResponse, error) {
	// Return success response
	return &apiv1.HealthCheckResponse{
		Status: "OK",
	}, nil
}
