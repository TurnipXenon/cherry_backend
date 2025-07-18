# Cherry Backend Testing Documentation

This document provides detailed information about testing the Cherry Backend service.

## Types of Tests

The Cherry Backend project includes two main types of tests:

### Unit Tests

Unit tests are located in the `internal/server` directory and test individual components of the server in isolation. These tests include:

- **Service Tests**: Tests for the service implementations (e.g., `TodoistServiceImpl`, `HealthServiceImpl`)
  - `todoist_service_test.go`: Tests the `ProcessWebhook` method of the `TodoistServiceImpl`
  - `health_service_test.go`: Tests the `Check` method of the `HealthServiceImpl`

- **Mock Tests**: Tests that use mock implementations of the services
  - `mock_test.go`: Tests using mock implementations of the services
  - `mock/mock_test.go`: Additional tests using mock implementations in a separate package

Unit tests focus on testing the business logic of individual components without dependencies on external systems. They are fast, reliable, and help ensure that each component works correctly in isolation.

### Integration/End-to-End Tests

Integration tests are located in the `test` directory and test the entire system, including the API endpoints and the interaction between components. These tests include:

- **API Client Tests**: Tests that use the API clients to send requests to the server
  - `client/client.go`: Tests the Todoist webhook and health check endpoints

- **Webhook Simulator**: A tool that simulates webhook events from Todoist
  - `webhook/simulate_webhook.go`: Sends simulated webhook events to the server

Integration tests focus on testing the system as a whole, including the HTTP endpoints, request/response handling, and the interaction between components. They help ensure that the entire system works correctly together.

### Key Differences

| Aspect | Unit Tests | Integration Tests |
|--------|------------|-------------------|
| **Location** | `internal/server` directory | `test` directory |
| **Focus** | Individual components | Entire system |
| **Dependencies** | Minimal (mocks used) | Real dependencies |
| **Speed** | Fast | Slower |
| **Scope** | Narrow (specific functions) | Broad (end-to-end workflows) |
| **When to use** | During development to test business logic | Before deployment to test the entire system |

## Comprehensive Testing Strategy

For thorough testing of your webhook implementation, it's recommended to use a combination of the tools provided in this project:

### Automated Testing with Test Runner Scripts

The project includes scripts that automate the process of running the server and test client together. These scripts handle starting the server, running the client against it, and optionally running the webhook simulator.

#### Windows (PowerShell)

```powershell
# From the project root
.\test\run_tests.ps1
```

#### Linux/macOS (Bash)

```bash
# From the project root
chmod +x test/run_tests.sh  # Make the script executable (first time only)
./test/run_tests.sh
```

#### Script Options

Both scripts support the following options:

**Windows (PowerShell):**
```powershell
.\test\run_tests.ps1 -ServerPort 9090 -ClientHost "localhost" -RunWebhookSimulator -WebhookEvent "item:completed" -Secret "your_secret"
```

**Linux/macOS (Bash):**
```bash
./test/run_tests.sh --port=9090 --host="localhost" --webhook --event="item:completed" --secret="your_secret"
```

| Option | Description | Default |
|--------|-------------|---------|
| ServerPort / --port | The port the server will listen on | 8080 |
| ClientHost / --host | The hostname the client will connect to | localhost |
| RunWebhookSimulator / --webhook | Whether to run the webhook simulator | false |
| WebhookEvent / --event | The event type for the webhook simulator | item:added |
| Secret / --secret | The secret for webhook signature verification | (empty) |

## Testing the API Clients

The project includes a test client in the `client` directory that demonstrates how to use the API clients. This is useful for testing the API endpoints without writing a full application.

### Running the Test Client

```bash
# Navigate to the client directory
cd client

# Build the test client
go build -o client client.go

# Run the test client
./client
```

Note: The test client expects the server to be running on localhost:8080. If your server is running on a different host or port, you'll need to modify the client.go file.

### What the Test Client Does

The test client performs two tests:

1. **Todoist Client Test**: Sends a simulated webhook event to the server and prints the response.
   - Creates a Todoist client with the server URL
   - Creates a webhook request with an "item:added" event
   - Sends the request to the server
   - Prints the success status and message from the response

2. **Health Client Test**: Performs a health check on the server and prints the response.
   - Creates a Health client with the server URL
   - Sends a health check request to the server
   - Prints the status from the response

### Modifying the Test Client

You can modify the test client to test different scenarios:

- Change the event type in the webhook request (e.g., "item:updated", "item:deleted", "item:completed")
- Change the user ID or event data in the webhook request
- Change the server URL to test a different server

Example modification to test an "item:completed" event:

```go
// Create a webhook request
request := &api.TodoistWebhookRequest{
    EventName: "item:completed",  // Changed from "item:added"
    UserID:    "123456",
    EventData: `{"item_id": "789", "completed_at": "2023-12-31T12:00:00Z"}`,  // Added completed_at
    Version:   "1.0",
}
```

## Testing the Webhook

### Using the Webhook Simulator

A webhook simulator is included in the `webhook` directory to help you test your webhook implementation without needing a real Todoist integration. This tool simulates Todoist webhook events by sending HTTP POST requests to your server with properly formatted payloads.

#### Building the Simulator

```bash
# Navigate to the webhook directory
cd webhook

# Build the simulator
go build -o simulate_webhook simulate_webhook.go
```

#### Running the Simulator

Basic usage with default settings:

```bash
# Run with default settings (sends an item:added event to localhost:8080)
./simulate_webhook
```

This will:
1. Create a webhook payload for an "item:added" event
2. Send it to http://localhost:8080/webhooks/todoist
3. Print the response from the server

#### Advanced Usage

You can customize the simulator's behavior using command-line flags:

```bash
# Run with custom parameters
./simulate_webhook -url="http://your-server.com/webhooks/todoist" -secret="your_secret" -event="item:completed" -user="67890"
```

Available parameters:
- `-url`: The URL to send the webhook to (default: "http://localhost:8080/webhooks/todoist")
- `-secret`: Your Todoist client secret for signature verification
- `-event`: Event type (item:added, item:updated, item:deleted, item:completed)
- `-user`: User ID for the webhook payload

#### How the Simulator Works

The simulator performs the following steps:

1. **Creates a sample task**: Generates a test task with a title, description, due date, and priority
2. **Builds a webhook payload**: Packages the task data into a webhook payload with the specified event type and user ID
3. **Calculates a signature** (if a secret is provided): Uses HMAC-SHA256 to create a signature that Todoist would use to verify the webhook
4. **Sends an HTTP POST request**: Submits the payload to your server with the appropriate headers
5. **Displays the response**: Shows the HTTP status code and a success message

#### Example Use Cases

1. **Testing signature verification**:
   ```bash
   ./simulate_webhook -secret="your_todoist_client_secret"
   ```
   This tests if your server correctly validates the webhook signature.

2. **Testing different event types**:
   ```bash
   ./simulate_webhook -event="item:completed"
   ./simulate_webhook -event="item:updated"
   ./simulate_webhook -event="item:deleted"
   ```
   These commands test how your server handles different types of Todoist events.

3. **Testing with a specific user ID**:
   ```bash
   ./simulate_webhook -user="12345" -event="item:added"
   ```
   This tests how your server processes events for a specific user.

4. **Testing a deployed server**:
   ```bash
   ./simulate_webhook -url="https://your-production-server.com/webhooks/todoist" -secret="your_production_secret"
   ```
   This tests your production server with a simulated webhook.

Note: The simulator uses the `models` package from the main application, so make sure you have the correct module structure set up.

### Using a Real Todoist Integration

To test with a real Todoist integration:

1. Make changes in your Todoist account that trigger the events you've subscribed to
2. Check your server logs for webhook reception and processing messages
3. Verify that the signature verification is working correctly

## Comprehensive Testing Strategy

For thorough testing of your webhook implementation, it's recommended to use a combination of the tools provided in this project:

### Automated Testing with Test Runner Scripts

The project includes scripts that automate the process of running the server and test client together. These scripts handle starting the server, running the client against it, and optionally running the webhook simulator.

#### Windows (PowerShell)

```powershell
# From the project root
.\test\run_tests.ps1
```

#### Linux/macOS (Bash)

```bash
# From the project root
chmod +x test/run_tests.sh  # Make the script executable (first time only)
./test/run_tests.sh
```

#### Script Options

Both scripts support the following options:

**Windows (PowerShell):**
```powershell
.\test\run_tests.ps1 -ServerPort 9090 -ClientHost "localhost" -RunWebhookSimulator -WebhookEvent "item:completed" -Secret "your_secret"
```

**Linux/macOS (Bash):**
```bash
./test/run_tests.sh --port=9090 --host="localhost" --webhook --event="item:completed" --secret="your_secret"
```

| Option | Description | Default |
|--------|-------------|---------|
| ServerPort / --port | The port the server will listen on | 8080 |
| ClientHost / --host | The hostname the client will connect to | localhost |
| RunWebhookSimulator / --webhook | Whether to run the webhook simulator | false |
| WebhookEvent / --event | The event type for the webhook simulator | item:added |
| Secret / --secret | The secret for webhook signature verification | (empty) |

### Manual End-to-End Testing Workflow

If you prefer to run the components manually instead of using the automated scripts above, follow these steps:

1. **Start the server**:
   ```bash
   # From the project root
   go run main.go
   ```

2. **Test the health endpoint** using the test client:
   ```bash
   # From the test directory
   cd client
   go build -o client client.go
   ./client
   ```
   Verify that the health check returns "OK".

3. **Test webhook processing** using the webhook simulator:
   ```bash
   # From the webhook directory
   go build -o simulate_webhook simulate_webhook.go
   ./simulate_webhook
   ```
   Verify that the server logs show the webhook being received and processed.

4. **Test signature verification** by providing a secret:
   ```bash
   # First, set the environment variable in your server
   export TODOIST_CLIENT_SECRET="your_test_secret"

   # Then restart the server and run the simulator with the same secret
   ./simulate_webhook -secret="your_test_secret"
   ```
   Verify that the server logs show the signature being verified successfully.

### Troubleshooting

If you encounter issues while testing:

1. **Server not starting**:
   - Check if the port is already in use
   - Verify that all required environment variables are set
   - Check the logs for any error messages

2. **Client test failing**:
   - Ensure the server is running on the expected host and port
   - Check network connectivity between the client and server
   - Verify that the API endpoints are correctly implemented

3. **Webhook simulator failing**:
   - Check that the URL is correct and the server is accessible
   - If using a secret, ensure it matches what the server expects
   - Verify that the payload format matches what the server expects

4. **Signature verification failing**:
   - Double-check that the secret is exactly the same on both sides
   - Ensure the environment variable is correctly set and loaded
   - Check the server logs for details about the signature mismatch

5. **Module not found errors**:
   - Run `go mod tidy` to ensure all dependencies are correctly resolved
   - Check that the module path in go.mod matches your import statements

6. **Automated test script errors**:
   - If you see an error like `package test/client/temp_client.go is not in GOROOT`, this is because Go is interpreting the file path as a package path. The test scripts have been updated to fix this issue.
   - If you're still encountering this error, make sure you're using the latest version of the test scripts.
   - For other script-related issues, check that you have the necessary permissions to execute the scripts and that all dependencies are installed.
