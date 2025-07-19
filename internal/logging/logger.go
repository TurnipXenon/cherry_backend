package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Warn(format string, args ...interface{})
}

// LoggerImpl implements the Logger interface
type LoggerImpl struct {
	mu       sync.Mutex
	file     *os.File
	basePath string
	date     string
}

// NewLogger creates a new logger instance
func NewLogger() (*LoggerImpl, error) {
	basePath := getLogPath()

	// Create the log directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logger := &LoggerImpl{
		basePath: basePath,
	}

	// Initialize the log file
	if err := logger.rotateLogFileIfNeeded(); err != nil {
		return nil, err
	}

	return logger, nil
}

// getLogPath returns the platform-specific path for log files
func getLogPath() string {
	// Check if the environment variable is set for testing
	if path := os.Getenv("CHERRY_LOG_PATH"); path != "" {
		return path
	}

	if runtime.GOOS == "windows" {
		// On Windows, store logs in a 'logs' directory in the repository
		return filepath.Join(".", "logs")
	} else {
		// On Linux/Unix, store logs in /var/log/cherry
		return "/var/log/cherry"
	}
}

// rotateLogFileIfNeeded checks if the log file needs to be rotated and does so if necessary
func (l *LoggerImpl) rotateLogFileIfNeeded() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	currentDate := time.Now().Format("2006-01-02")

	// If the date has changed or the file is not open, rotate the log file
	if l.date != currentDate || l.file == nil {
		// Close the current file if it's open
		if l.file != nil {
			l.file.Close()
			l.file = nil
		}

		// Open a new log file for the current date
		logFilePath := filepath.Join(l.basePath, fmt.Sprintf("cherry-%s.log", currentDate))
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		l.file = file
		l.date = currentDate
	}

	return nil
}

// log writes a log message to the log file
func (l *LoggerImpl) log(level, format string, args ...interface{}) {
	if err := l.rotateLogFileIfNeeded(); err != nil {
		fmt.Fprintf(os.Stderr, "Error rotating log file: %v\n", err)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Write to the log file
	if _, err := l.file.WriteString(logLine); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to log file: %v\n", err)
	}

	// Also write to stdout for convenience
	fmt.Print(logLine)
}

// Info logs an informational message
func (l *LoggerImpl) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Error logs an error message
func (l *LoggerImpl) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

// Debug logs a debug message
func (l *LoggerImpl) Debug(format string, args ...interface{}) {
	l.log("DEBUG", format, args...)
}

// Warn logs a warning message
func (l *LoggerImpl) Warn(format string, args ...interface{}) {
	l.log("WARN", format, args...)
}

// Close closes the log file
func (l *LoggerImpl) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}

	return nil
}
