package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/unwale/skingen/pkg/logging"
)

func LoggingMiddleware(baseLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := uuid.New().String()
			r.Header.Set("X-Correlation-ID", correlationID)

			logger := baseLogger.With(
				slog.String("correlation_id", correlationID),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)
			ctx := logging.WithLogger(r.Context(), logger)
			r = r.WithContext(ctx)

			logger.Info("Starting request")
			start := time.Now()

			defer func() {
				duration := time.Since(start)

				if err := recover(); err != nil {
					logger.Error("handler panicked",
						"error", err,
						slog.Int("status", http.StatusInternalServerError),
						slog.Duration("duration", duration),
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				} else {
					logger.Info("request completed",
						slog.Duration("duration", duration),
					)
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}
