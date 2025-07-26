package interceptors

// In services/task-service/internal/api/interceptors.go (or similar)

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/unwale/skingen/pkg/constants"
	"github.com/unwale/skingen/pkg/contextutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingInterceptor(baseLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var correlationID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(constants.CorrelationIDKey); len(vals) > 0 {
				correlationID = vals[0]
			}
		}

		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		logger := baseLogger.With(
			slog.String("correlation_id", correlationID),
			slog.String("grpc.method", info.FullMethod),
		)

		ctx = contextutil.WithLogger(ctx, logger)
		ctx = contextutil.WithCorrelationID(ctx, correlationID)

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			logger.Error("gRPC call failed",
				"error", err,
				slog.Duration("duration", duration),
			)
		} else {
			logger.Info("gRPC call completed",
				slog.Duration("duration", duration),
			)
		}

		return resp, err
	}
}
