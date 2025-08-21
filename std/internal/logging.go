package internal

import (
	"log/slog"
	"os"
)

// TODO: integrate with otel

func CreateLogger(name string) *slog.Logger {
	return slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug})).With(slog.String("logger_name", name))
}
