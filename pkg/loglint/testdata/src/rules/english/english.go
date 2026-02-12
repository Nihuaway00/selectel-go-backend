package english

import (
	"log/slog"

	"go.uber.org/zap"
)

func example() {
	slog.Info("запуск сервера") // want "log message should contain only English letters"

	var logger *slog.Logger
	logger.Error("ошибка подключения") // want "log message should contain only English letters"

	var zapLogger *zap.Logger
	zapLogger.Info("не удалось подключиться") // want "log message should contain only English letters"
}

func edgeCases() {
	slog.Info("ASCII only 123") // ok
	slog.Info("café")           // want "log message should contain only English letters"
	slog.Info("über")           // want "log message should contain only English letters"
}
