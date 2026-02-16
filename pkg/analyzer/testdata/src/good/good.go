package good

import (
	"log/slog"
)

func testSlog() {
	slog.Info("starting server")
	slog.Error("failed to connect to database")
	slog.Debug("processing completed")
	slog.Warn("timeout occurred")
}
