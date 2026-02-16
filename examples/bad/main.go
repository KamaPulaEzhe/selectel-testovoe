package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

func main() {
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	zapLogger := zap.NewExample()

	password := "supersecret"
	apiKey := "12345"

	slogLogger.Info("Server started on port 8080")
	slogLogger.Error("Failed to connect to Database")
	slogLogger.Info("привет мир")
	slogLogger.Warn("warning!!!")
	slogLogger.Debug("user password: " + password)
	slogLogger.Info("api_key: " + apiKey)

	zapLogger.Info("Server started")
	zapLogger.Error("ошибка подключения")
	zapLogger.Warn("loading...")
	zapLogger.Debug("token: abc123")
}
