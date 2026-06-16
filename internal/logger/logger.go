package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Log is the global structured logger instance used across the application
var Log *slog.Logger

// InitLogger configures and initializes the structured JSON logging system
func InitLogger() error {
	logDir := "logs"
	logFile := filepath.Join(logDir, "out.log")

	// Ensure the log directory exists with correct permissions
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Open the log file in write-only, append mode, or create it if it doesn't exist
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Write logs to both standard system output and the physical log file simultaneously
	multiWriter := io.MultiWriter(os.Stdout, file)

	// Configure the JSON handler options to format logs as valid structured JSON
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Initialize the global logger and set it as the application default
	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}