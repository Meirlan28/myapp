package logger

import (
	"log/slog"
	"os"
)

func New(cfg Config) *slog.Logger {
	var handler slog.Handler
	var level slog.Level
	switch cfg.Level {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}
	handler = slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: level,
		},
	)

	return slog.New(handler)
}
