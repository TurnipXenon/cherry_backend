# Cherry Backend Test Runner
#
# This script runs the Cherry Backend server and test client concurrently.
# It handles setting up the correct ports and ensures proper cleanup.

param (
    [int]$ServerPort = 8080,
    [string]$ClientHost = "localhost",
    [switch]$RunWebhookSimulator = $false,
    [string]$WebhookEvent = "item:added",
    [string]$Secret = ""
)

# Display script banner
Write-Host "========================================"
Write-Host "Cherry Backend Test Runner"
Write-Host "========================================"
Write-Host "Server Port: $ServerPort"
Write-Host "Client Host: $ClientHost"
Write-Host "Run Webhook Simulator: $RunWebhookSimulator"
if ($RunWebhookSimulator) {
    Write-Host "Webhook Event: $WebhookEvent"
    if ($Secret) {
        Write-Host "Secret: [PROVIDED]"
    } else {
        Write-Host "Secret: [NONE]"
    }
}
Write-Host "========================================"
Write-Host ""

# Create a temporary client file with the correct server URL
$clientCode = Get-Content -Path "test\client\client.go" -Raw
$serverUrl = "http://${ClientHost}:${ServerPort}"
$modifiedClientCode = $clientCode -replace 'http://localhost:8080', $serverUrl

# Save the modified client code to a temporary file
$tempClientPath = "test\client\temp_client.go"
Set-Content -Path $tempClientPath -Value $modifiedClientCode

# Build the temporary client
Write-Host "Building test client..."
Push-Location "test\client"
# Use the correct syntax for building a specific Go file
go build -o temp_client.exe temp_client.go
if (-not $?) {
    Write-Host "Failed to build test client" -ForegroundColor Red
    Remove-Item -Path $tempClientPath -ErrorAction SilentlyContinue
    Pop-Location
    exit 1
}
Pop-Location

# Start the server in a background job
Write-Host "Starting server on port $ServerPort..."
$serverJob = Start-Job -ScriptBlock {
    param($port)
    $env:PORT = $port
    $env:TODOIST_CLIENT_SECRET = "test_secret"
    Set-Location $using:PWD
    go run main.go
} -ArgumentList $ServerPort

# Wait for the server to start
Write-Host "Waiting for server to start..."
Start-Sleep -Seconds 3

# Run the client
Write-Host "Running test client..."
Push-Location "test\client"
.\temp_client.exe
$clientExitCode = $LASTEXITCODE
Pop-Location

# Run the webhook simulator if requested
if ($RunWebhookSimulator) {
    Write-Host "Running webhook simulator..."
    Push-Location "test\webhook"

    # Build the webhook simulator if it doesn't exist
    if (-not (Test-Path "simulate_webhook.exe")) {
        Write-Host "Building webhook simulator..."
        go build -o simulate_webhook.exe simulate_webhook.go
        if (-not $?) {
            Write-Host "Failed to build webhook simulator" -ForegroundColor Red
            Pop-Location
            # Clean up
            Remove-Item -Path "test\client\temp_client.exe" -ErrorAction SilentlyContinue
            Remove-Item -Path "test\client\temp_client.go" -ErrorAction SilentlyContinue
            Stop-Job -Job $serverJob
            Remove-Job -Job $serverJob
            exit 1
        }
    }

    # Run the webhook simulator
    $webhookUrl = "http://${ClientHost}:${ServerPort}/webhooks/todoist"
    if ($Secret) {
        .\simulate_webhook.exe -url=$webhookUrl -event=$WebhookEvent -secret=$Secret
    } else {
        .\simulate_webhook.exe -url=$webhookUrl -event=$WebhookEvent
    }
    Pop-Location
}

# Clean up
Write-Host "Cleaning up..."
Remove-Item -Path "test\client\temp_client.exe" -ErrorAction SilentlyContinue
Remove-Item -Path "test\client\temp_client.go" -ErrorAction SilentlyContinue

# Stop the server
Write-Host "Stopping server..."
Stop-Job -Job $serverJob
Remove-Job -Job $serverJob

# Display summary
Write-Host ""
Write-Host "========================================"
Write-Host "Test Summary"
Write-Host "========================================"
if ($clientExitCode -eq 0) {
    Write-Host "Client tests: SUCCESS" -ForegroundColor Green
} else {
    Write-Host "Client tests: FAILED (Exit code: $clientExitCode)" -ForegroundColor Red
}
Write-Host "========================================"

exit $clientExitCode
