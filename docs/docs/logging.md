# Logging System

Cherry Backend includes a simple but effective logging system that writes logs to files grouped by day. This document explains how the logging system works and how to use it in your code.

## Overview

The logging system is designed to:

- Write logs to files that are grouped by day (one file per day)
- Store logs in platform-specific locations:
  - On Windows: in a `logs` directory in the repository
  - On Linux: in `/var/log/cherry`
- Provide basic logging levels: Info, Error, Debug, and Warn
- Be thread-safe for concurrent logging
- Automatically rotate log files based on the date

## Log File Format

Log files are named using the format `cherry-YYYY-MM-DD.log`, where `YYYY-MM-DD` is the date. For example, logs for July 18, 2025 would be stored in `cherry-2025-07-18.log`.

Each log entry has the following format:

```
[YYYY-MM-DD HH:MM:SS] [LEVEL] MESSAGE
```

For example:

```
[2025-07-18 15:04:05] [INFO] Processing webhook event: item:added
[2025-07-18 15:04:05] [ERROR] Failed to process webhook: invalid signature
```

## Using the Logger

### Creating a Logger

To use the logging system in your code, you need to create a logger instance:

```go
import "cherry_backend/internal/logging"

// Create a new logger
logger, err := logging.NewLogger()
if err != nil {
    // Handle error
}
defer logger.Close() // Don't forget to close the logger when done
```

### Logging Messages

The logger provides four methods for logging messages at different levels:

```go
// Informational messages
logger.Info("Processing webhook event: %s", eventName)

// Error messages
logger.Error("Failed to process webhook: %v", err)

// Debug messages
logger.Debug("Request details: %+v", request)

// Warning messages
logger.Warn("TODOIST_CLIENT_SECRET not set")
```

### In the TodoistServiceImpl

The `TodoistServiceImpl` already includes a logger that is initialized in the constructor:

```go
// Create a new TodoistServiceImpl with a logger
service, err := NewTodoistServiceImpl()
if err != nil {
    // Handle error
}

// Use the logger
service.Logger.Info("Processing webhook event: %s", eventName)
```

## Log File Locations

### Windows

On Windows, log files are stored in a `logs` directory in the repository root:

```
C:\path\to\cherry_backend\logs\cherry-2025-07-18.log
```

### Linux

On Linux, log files are stored in `/var/log/cherry`:

```
/var/log/cherry/cherry-2025-07-18.log
```

Make sure the application has write permissions to this directory.

## Implementation Details

The logging system is implemented in the `internal/logging` package and consists of:

- A `Logger` interface that defines the logging methods
- A `LoggerImpl` struct that implements the interface
- Functions for creating a new logger, rotating log files, and writing log messages

The implementation uses a mutex to ensure thread safety and checks the date on each log write to determine if a log rotation is needed.

## Limitations and Future Improvements

The current logging system is simple and meets the basic requirements, but it has some limitations:

1. **No log level filtering**: All log messages are written to the log file regardless of their level. A future improvement could add the ability to set a minimum log level.

2. **No log compression**: Old log files are not compressed or archived. A future improvement could add automatic compression of old log files.

3. **No log retention policy**: Log files are never deleted. A future improvement could add a retention policy to delete or archive old log files.

4. **Limited configuration**: The log file location and format are hardcoded. A future improvement could make these configurable.

5. **No structured logging**: The current system uses simple text logging. A future improvement could add support for structured logging formats like JSON.

6. **No remote logging**: Logs are only written to local files. A future improvement could add support for sending logs to remote systems.

7. **No log correlation**: There's no way to correlate logs across different services or requests. A future improvement could add support for trace IDs or correlation IDs.

## Testing

The logging system includes tests that verify:

1. A logger can be created
2. Log messages are written to the log file
3. Log files are rotated based on the date

To run the tests:

```bash
go test -v ./internal/logging
```

The tests use a temporary directory for log files to avoid interfering with the actual application logs.