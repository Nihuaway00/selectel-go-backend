package literal

import "log/slog"

func example() {
	msg := "dynamic"
	slog.Info(msg) // want "log message should be a string literal"
}

const constMsg = "server" + " started"

func edgeCases() {
	slog.Info(constMsg) // ok: constant string

	dynamic := "dynamic"
	slog.Info(dynamic) // want "log message should be a string literal"
}
