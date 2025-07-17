# Cherry Backend

A backend service for the Cherry project.

## Description

Cherry Backend is a Go-based service that provides API endpoints for the Cherry application.

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Git
- Protocol Buffers compiler (protoc)
- Make (for Windows users, see installation instructions below)

#### Installing Make on Windows

If you're using Windows, you can install Make using Chocolatey:

1. Install Chocolatey (if not already installed)
   ```
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
   ```

2. Install Make
   ```
   choco install make
   ```

3. Verify installation
   ```
   make --version
   ```

#### Installing Protocol Buffers Compiler

1. Install Protocol Buffers compiler for your operating system:
   - For Windows (using Chocolatey):
     ```
     choco install protoc
     ```
   - For macOS (using Homebrew):
     ```
     brew install protobuf
     ```
   - For Linux:
     ```
     apt-get install -y protobuf-compiler
     ```

2. Install Go plugins for Protocol Buffers:
   ```
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. Verify installation:
   ```
   protoc --version
   ```

### Installation

1. Clone the repository
   ```
   git clone https://github.com/yourusername/cherry_backend.git
   ```

2. Navigate to the project directory
   ```
   cd cherry_backend
   ```

3. Set up the configuration
   ```
   cp configs/local.env.example configs/local.env
   # Edit configs/local.env with your settings
   ```

4. Build the project
   ```
   go build
   ```

5. Run the application
   ```
   ./cherry_backend
   ```

## Usage

### Running the Service

```bash
# Build the application
go build -o cherry_backend

# Run the application
./cherry_backend
```

The server will start on the configured port (default: 8080) and will be ready to receive webhook notifications from Todoist.

## API Documentation

The API is defined using Protocol Buffers (protobuf). The protobuf definitions can be found in the `proto` directory.

### Generating Clients from Protobuf

The project uses Protocol Buffers to define the API and generate client code. To generate the Go code from the protobuf definitions:

1. Make sure you have installed the prerequisites (protoc and Go plugins)
2. Run the following command from the project root:
   ```
   make proto
   ```

This will generate the Go code in the `pkg/api/v1` directory. The generated code includes:
- Message definitions (e.g., `TodoistWebhookRequest`, `TodoistWebhookResponse`)
- Service interfaces (e.g., `TodoistServiceClient`, `HealthServiceClient`)
- gRPC client implementations

If you want to clean the generated files, you can run:
```
make clean
```

#### Using the Makefile with Git Bash

If you're using Git Bash on Windows, the Makefile will automatically detect it and use Unix-like commands. If the automatic detection doesn't work, you can manually force the use of Unix-like commands by setting the `USE_UNIX_COMMANDS` variable:

```
make proto USE_UNIX_COMMANDS=1
```

or

```
make clean USE_UNIX_COMMANDS=1
```

Note: The REST clients in `pkg/api` directory are manually implemented and use the generated message types. If you make changes to the protobuf definitions, you may need to update these clients accordingly.

### Webhook Endpoints

#### Todoist Webhook

- **URL**: `/webhooks/todoist`
- **Method**: `POST`
- **Description**: Receives webhook notifications from Todoist when events occur in a user's Todoist account.
- **Headers**:
  - `X-Todoist-Hmac-SHA256`: HMAC-SHA256 signature for request verification

The webhook handler supports various Todoist event types, including:
- `item:added`
- `item:updated`
- `item:deleted`
- `item:completed`

### Health Check

- **URL**: `/health`
- **Method**: `GET`
- **Description**: Simple health check endpoint that returns a status message if the service is running.

## Client Generator

The project includes a client generator that creates REST clients for the API. The client generator can be found in the `pkg/api` directory.

### Usage

```go
// Create a Todoist client
todoistClient := api.NewTodoistClient("http://localhost:8080")

// Create a webhook request
request := &api.TodoistWebhookRequest{
    EventName: "item:added",
    UserID:    "123456",
    EventData: `{"item_id": "789"}`,
    Version:   "1.0",
}

// Process the webhook
response, err := todoistClient.ProcessWebhook(request)
if err != nil {
    log.Fatalf("Failed to process webhook: %v", err)
}

// Print the response
fmt.Printf("Success: %v, Message: %s\n", response.Success, response.Message)
```

```go
// Create a health check client
healthClient := api.NewHealthClient("http://localhost:8080")

// Perform a health check
response, err := healthClient.Check()
if err != nil {
    log.Fatalf("Failed to perform health check: %v", err)
}

// Print the response
fmt.Printf("Status: %s\n", response.Status)
```

## Todoist Webhook Integration

### Setup Instructions

1. **Configure Environment Variables**:
   - Copy `.env.example` to `.env`
   - Set your Todoist client secret in the `.env` file:
     ```
     TODOIST_CLIENT_SECRET=your_actual_client_secret
     ```

2. **Register the Webhook in Todoist**:
   - Go to the [Todoist Developer Console](https://developer.todoist.com/appconsole.html)
   - Create a new app or select an existing one
   - In the app settings, add a new webhook with the following details:
     - **Webhook URL**: `https://your-server-domain.com/webhooks/todoist`
     - **Events**: Select the events you want to receive notifications for

3. **Ensure Your Server is Publicly Accessible**:
   - Todoist needs to be able to reach your server to send webhook notifications
   - You may need to set up port forwarding, use a service like ngrok for development, or deploy to a cloud provider

### Testing the Webhook

#### Using the Webhook Simulator

A webhook simulator is included in the `test` directory to help you test your webhook implementation without needing a real Todoist integration:

```bash
# Build the simulator
cd test
go build -o simulate_webhook simulate_webhook.go

# Run the simulator with default settings (sends an item:added event to localhost:8080)
./simulate_webhook

# Run with custom parameters
./simulate_webhook -url="http://your-server.com/webhooks/todoist" -secret="your_secret" -event="item:completed" -user="67890"
```

Note: The simulator uses the `models` package from the main application, so make sure you have the correct module structure set up.

Available parameters:
- `-url`: The URL to send the webhook to (default: "http://localhost:8080/webhooks/todoist")
- `-secret`: Your Todoist client secret for signature verification
- `-event`: Event type (item:added, item:updated, item:deleted, item:completed)
- `-user`: User ID for the webhook payload

#### Using a Real Todoist Integration

To test with a real Todoist integration:

1. Make changes in your Todoist account that trigger the events you've subscribed to
2. Check your server logs for webhook reception and processing messages
3. Verify that the signature verification is working correctly

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- List any acknowledgments here
