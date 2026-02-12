package special

import (
	"log/slog"

	"go.uber.org/zap"
)

func example() {
	slog.Info("server started!")                  // want "log message should not contain special characters or emoji"
	slog.Info("warning: something went wrong...") // want "log message should not contain special characters or emoji"
	slog.Info("server started 🚀")                 // want "log message should not contain special characters or emoji"

	var zapLogger *zap.Logger
	zapLogger.Error("connection failed!!!") // want "log message should not contain special characters or emoji"
}

func edgeCases() {
	slog.Info("cache miss: key")  // want "log message should not contain special characters or emoji"
	slog.Info("hello_world")      // want "log message should not contain special characters or emoji"
	slog.Info("tab\tseparated")   // want "log message should not contain special characters or emoji"
	slog.Info("newline\nmessage") // want "log message should not contain special characters or emoji"
	slog.Info("letters 123 only") // ok
}
