package models

// TodoistWebhookPayload represents the structure of a Todoist webhook payload
type TodoistWebhookPayload struct {
	EventName string                 `json:"event_name"`
	UserID    string                 `json:"user_id"`
	EventData map[string]interface{} `json:"event_data"`
	Version   string                 `json:"version"`
}