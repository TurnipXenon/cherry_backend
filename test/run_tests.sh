#!/bin/bash
# Cherry Backend Test Runner
#
# This script runs the Cherry Backend server and test client concurrently.
# It handles setting up the correct ports and ensures proper cleanup.

# Default values
SERVER_PORT=8080
CLIENT_HOST="localhost"
RUN_WEBHOOK_SIMULATOR=false
WEBHOOK_EVENT="item:added"
SECRET=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --port=*)
      SERVER_PORT="${1#*=}"
      shift
      ;;
    --host=*)
      CLIENT_HOST="${1#*=}"
      shift
      ;;
    --webhook)
      RUN_WEBHOOK_SIMULATOR=true
      shift
      ;;
    --event=*)
      WEBHOOK_EVENT="${1#*=}"
      shift
      ;;
    --secret=*)
      SECRET="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Display script banner
echo "========================================"
echo "Cherry Backend Test Runner"
echo "========================================"
echo "Server Port: $SERVER_PORT"
echo "Client Host: $CLIENT_HOST"
echo "Run Webhook Simulator: $RUN_WEBHOOK_SIMULATOR"
if [ "$RUN_WEBHOOK_SIMULATOR" = true ]; then
    echo "Webhook Event: $WEBHOOK_EVENT"
    if [ -n "$SECRET" ]; then
        echo "Secret: [PROVIDED]"
    else
        echo "Secret: [NONE]"
    fi
fi
echo "========================================"
echo ""

# Create a temporary client file with the correct server URL
CLIENT_CODE=$(cat test/client/client.go)
SERVER_URL="http://${CLIENT_HOST}:${SERVER_PORT}"
MODIFIED_CLIENT_CODE=${CLIENT_CODE//http:\/\/localhost:8080/$SERVER_URL}

# Save the modified client code to a temporary file
TEMP_CLIENT_PATH="test/client/temp_client.go"
echo "$MODIFIED_CLIENT_CODE" > "$TEMP_CLIENT_PATH"

# Build the temporary client
echo "Building test client..."
pushd test/client > /dev/null
# Use the correct syntax for building a specific Go file
go build -o temp_client temp_client.go
if [ $? -ne 0 ]; then
    echo "Failed to build test client"
    rm -f "$TEMP_CLIENT_PATH"
    popd > /dev/null
    exit 1
fi
popd > /dev/null

# Start the server in the background
echo "Starting server on port $SERVER_PORT..."
export PORT=$SERVER_PORT
export TODOIST_CLIENT_SECRET="test_secret"
go run main.go &
SERVER_PID=$!

# Wait for the server to start
echo "Waiting for server to start..."
sleep 3

# Run the client
echo "Running test client..."
pushd test/client > /dev/null
./temp_client
CLIENT_EXIT_CODE=$?
popd > /dev/null

# Run the webhook simulator if requested
if [ "$RUN_WEBHOOK_SIMULATOR" = true ]; then
    echo "Running webhook simulator..."
    pushd test/webhook > /dev/null

    # Build the webhook simulator if it doesn't exist
    if [ ! -f "simulate_webhook" ]; then
        echo "Building webhook simulator..."
        go build -o simulate_webhook simulate_webhook.go
        if [ $? -ne 0 ]; then
            echo "Failed to build webhook simulator"
            popd > /dev/null
            # Clean up
            rm -f "test/client/temp_client"
            rm -f "$TEMP_CLIENT_PATH"
            kill $SERVER_PID
            exit 1
        fi
    fi

    # Run the webhook simulator
    WEBHOOK_URL="http://${CLIENT_HOST}:${SERVER_PORT}/webhooks/todoist"
    if [ -n "$SECRET" ]; then
        ./simulate_webhook -url="$WEBHOOK_URL" -event="$WEBHOOK_EVENT" -secret="$SECRET"
    else
        ./simulate_webhook -url="$WEBHOOK_URL" -event="$WEBHOOK_EVENT"
    fi
    popd > /dev/null
fi

# Clean up
echo "Cleaning up..."
rm -f "test/client/temp_client"
rm -f "$TEMP_CLIENT_PATH"

# Stop the server
echo "Stopping server..."
kill $SERVER_PID

# Display summary
echo ""
echo "========================================"
echo "Test Summary"
echo "========================================"
if [ $CLIENT_EXIT_CODE -eq 0 ]; then
    echo "Client tests: SUCCESS"
else
    echo "Client tests: FAILED (Exit code: $CLIENT_EXIT_CODE)"
fi
echo "========================================"

exit $CLIENT_EXIT_CODE
