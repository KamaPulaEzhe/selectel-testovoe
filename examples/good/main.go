package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

func main() {
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	zapLogger := zap.NewExample()

	userID := 12345

	slogLogger.Info("starting server on port 8080")
	slogLogger.Error("failed to connect to database")
	slogLogger.Debug("processing request", "user_id", userID)
	slogLogger.Warn("timeout occurred")

	zapLogger.Info("request completed")
	zapLogger.Error("connection refused")
	zapLogger.Debug("cache miss")
}
