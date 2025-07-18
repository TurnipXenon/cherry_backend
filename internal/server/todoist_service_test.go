package server

import (
	"context"
	"testing"

	apiv1 "cherry_backend/pkg/api/v1"
)

func TestProcessWebhook(t *testing.T) {
	// Create a new TodoistServiceImpl
	service := &TodoistServiceImpl{}

	// Create a context for the tests
	ctx := context.Background()

	// Define test cases
	testCases := []struct {
		name        string
		request     *apiv1.TodoistWebhookRequest
		wantSuccess bool
		wantMessage string
	}{
		{
			name: "item:added event",
			request: &apiv1.TodoistWebhookRequest{
				EventName: "item:added",
				UserId:    "123456",
				EventData: `{"item_id": "789"}`,
				Version:   "1.0",
			},
			wantSuccess: true,
			wantMessage: "Webhook received",
		},
		{
			name: "item:updated event",
			request: &apiv1.TodoistWebhookRequest{
				EventName: "item:updated",
				UserId:    "123456",
				EventData: `{"item_id": "789"}`,
				Version:   "1.0",
			},
			wantSuccess: true,
			wantMessage: "Webhook received",
		},
		{
			name: "item:deleted event",
			request: &apiv1.TodoistWebhookRequest{
				EventName: "item:deleted",
				UserId:    "123456",
				EventData: `{"item_id": "789"}`,
				Version:   "1.0",
			},
			wantSuccess: true,
			wantMessage: "Webhook received",
		},
		{
			name: "item:completed event",
			request: &apiv1.TodoistWebhookRequest{
				EventName: "item:completed",
				UserId:    "123456",
				EventData: `{"item_id": "789"}`,
				Version:   "1.0",
			},
			wantSuccess: true,
			wantMessage: "Webhook received",
		},
		{
			name: "unknown event",
			request: &apiv1.TodoistWebhookRequest{
				EventName: "unknown:event",
				UserId:    "123456",
				EventData: `{"item_id": "789"}`,
				Version:   "1.0",
			},
			wantSuccess: true,
			wantMessage: "Webhook received",
		},
	}

	// Run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the ProcessWebhook method
			response, err := service.ProcessWebhook(ctx, tc.request)

			// Check for errors
			if err != nil {
				t.Fatalf("ProcessWebhook returned an error: %v", err)
			}

			// Check the response
			if response.Success != tc.wantSuccess {
				t.Errorf("ProcessWebhook() success = %v, want %v", response.Success, tc.wantSuccess)
			}
			if response.Message != tc.wantMessage {
				t.Errorf("ProcessWebhook() message = %s, want %s", response.Message, tc.wantMessage)
			}
		})
	}
}