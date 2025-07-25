package contextutil

import (
	"context"
	"log/slog"
)

type contextKey string

var (
	loggerKey        = contextKey("logger")
	correlationIDKey = contextKey("correlationID")
)

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContextOrDefault(ctx context.Context, defaultLogger *slog.Logger) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return defaultLogger
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	if correlationID, ok := ctx.Value(correlationIDKey).(string); ok {
		return correlationID
	}
	return ""
}
