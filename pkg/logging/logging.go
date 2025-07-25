package logging

import (
	"context"
	"log/slog"
	"os"
)

type contextKey struct{}

var loggerKey = contextKey{}

func NewLogger(serviceName string, level string) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})
	logger := slog.New(handler)

	logger = logger.With(slog.String("service", serviceName))

	return logger
}

func FromContextOrDefault(ctx context.Context, defaultLogger *slog.Logger) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return defaultLogger
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
