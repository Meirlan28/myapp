package logger

import (
	"context"
	"log/slog"
	"os"
)

type loggerContextKey struct{}

var (
	key = loggerContextKey{}
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

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		logger,
	)
}

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(key).(*slog.Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}
