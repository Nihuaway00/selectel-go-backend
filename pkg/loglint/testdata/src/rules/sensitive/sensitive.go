package sensitive

import (
	"log/slog"

	"go.uber.org/zap"
)

func example() {
	slog.Info("user password: 123") // want "log message should not contain sensitive data"
	slog.Info("api_key=abcd")       // want "log message should not contain sensitive data"
	slog.Info("token: xyz")         // want "log message should not contain sensitive data"
	slog.Info("token validated")
	slog.Info("secret: value") // want "log message should not contain sensitive data"

	var zapLogger *zap.Logger
	zapLogger.Info("password = secret") // want "log message should not contain sensitive data"
}

func edgeCases() {
	slog.Info("token validated")     // ok
	slog.Info("token : value")       // want "log message should not contain sensitive data"
	slog.Info("password=secret")     // want "log message should not contain sensitive data"
	slog.Info("password123: secret") // ok
	slog.Info("api_key = secret")    // want "log message should not contain sensitive data"
	slog.Info("tokenized value")     // ok
	slog.Info("SECRET = value")      // want "log message should not contain sensitive data"
}
