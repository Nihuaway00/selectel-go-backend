package main

import "log/slog"

func main() {
	slog.Info("Starting server")             // should start with lowercase
	slog.Info("запуск сервера")              // non-English letters
	slog.Info("server started!")             // special chars
	slog.Info("user password: 123")          // sensitive data
	slog.Info("user secret=123")             // sensitive data
	slog.Info("server started successfully") // ok
}
