package lowercase

import (
	"log/slog"

	"go.uber.org/zap"
)

func example() {
	slog.Info("Starting server") // want "log message should start with a lowercase letter"
	slog.Info("starting server")

	var logger *slog.Logger
	logger.Error("Failed to connect") // want "log message should start with a lowercase letter"

	var zapLogger *zap.Logger
	zapLogger.Info("Hello") // want "log message should start with a lowercase letter"

	var sugar *zap.SugaredLogger
	sugar.Info("Hello") // want "log message should start with a lowercase letter"
}

func edgeCases() {
	slog.Info("  ok with leading spaces")
	slog.Info("")    // ok: empty
	slog.Info("   ") // ok: only spaces

	slog.Info("404 error") // want "log message should start with a lowercase letter"
	slog.Info("-failed")   // want "log message should start with a lowercase letter"
	slog.Info("🚀 launch")  // want "log message should start with a lowercase letter"
}
