package server

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	apiv1 "cherry_backend/pkg/api/v1"
)

// TestProcessWebhook tests the ProcessWebhook function with logging
func TestProcessWebhook(t *testing.T) {
	// Setup test environment
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create a new TodoistServiceImpl with a logger
	service, err := NewTodoistServiceImpl()
	if err != nil {
		t.Fatalf("Failed to create TodoistServiceImpl: %v", err)
	}

	// Create a test request
	request := &apiv1.TodoistWebhookRequest{
		EventName: "item:added",
		UserId:    "test-user",
		EventData: `{"item_id": "123", "content": "Test Item"}`,
		Version:   "1.0",
	}

	// Process the webhook
	response, err := service.ProcessWebhook(context.Background(), request)
	if err != nil {
		t.Fatalf("ProcessWebhook failed: %v", err)
	}

	// Check the response
	if !response.Success {
		t.Errorf("Expected success=true, got success=%v", response.Success)
	}
	if response.Message != "Webhook received" {
		t.Errorf("Expected message='Webhook received', got message='%s'", response.Message)
	}

	// Check that the log file contains the expected messages
	logContent, err := readLogFile(t)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check for expected log messages
	expectedMessages := []string{
		"Processing webhook event: item:added",
		"Item added by user test-user",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(logContent, msg) {
			t.Errorf("Log file does not contain expected message: %s", msg)
		}
	}
}

// TestProcessWebhookUnhandledEvent tests the ProcessWebhook function with an unhandled event type
func TestProcessWebhookUnhandledEvent(t *testing.T) {
	// Setup test environment
	setupTestEnv(t)
	defer cleanupTestEnv(t)

	// Create a new TodoistServiceImpl with a logger
	service, err := NewTodoistServiceImpl()
	if err != nil {
		t.Fatalf("Failed to create TodoistServiceImpl: %v", err)
	}

	// Create a test request with an unhandled event type
	request := &apiv1.TodoistWebhookRequest{
		EventName: "unknown:event",
		UserId:    "test-user",
		EventData: `{"item_id": "123", "content": "Test Item"}`,
		Version:   "1.0",
	}

	// Process the webhook
	response, err := service.ProcessWebhook(context.Background(), request)
	if err != nil {
		t.Fatalf("ProcessWebhook failed: %v", err)
	}

	// Check the response
	if !response.Success {
		t.Errorf("Expected success=true, got success=%v", response.Success)
	}

	// Check that the log file contains the expected messages
	logContent, err := readLogFile(t)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Check for expected log messages
	expectedMessage := "Unhandled event type: unknown:event"
	if !strings.Contains(logContent, expectedMessage) {
		t.Errorf("Log file does not contain expected message: %s", expectedMessage)
	}
}

// setupTestEnv sets up the test environment
func setupTestEnv(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := filepath.Join(os.TempDir(), "cherry_test_logs")
	os.MkdirAll(tempDir, 0755)
	t.Setenv("CHERRY_LOG_PATH", tempDir)
}

// cleanupTestEnv cleans up the test environment
func cleanupTestEnv(t *testing.T) {
	// Clean up the test log directory
	tempDir := filepath.Join(os.TempDir(), "cherry_test_logs")
	os.RemoveAll(tempDir)
}

// readLogFile reads the log file for the current date
func readLogFile(t *testing.T) (string, error) {
	// Get the log file path
	tempDir := filepath.Join(os.TempDir(), "cherry_test_logs")
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", err
	}

	// Find the log file (there should be only one in the test environment)
	if len(files) == 0 {
		return "", nil
	}

	// Read the log file
	logFilePath := filepath.Join(tempDir, files[0].Name())
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
