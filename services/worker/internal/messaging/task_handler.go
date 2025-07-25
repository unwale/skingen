package messaging

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/unwale/skingen/pkg/contextutil"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/worker/internal/core"
)

func CreateTaskCommandHandler(service core.WorkerService, baseLogger *slog.Logger) messaging.MessageHandler {
	return func(msg amqp091.Delivery) error {
		logger := baseLogger.With("correlation_id", msg.CorrelationId)
		ctx := contextutil.WithLogger(context.Background(), logger)
		ctx = contextutil.WithCorrelationID(ctx, msg.CorrelationId)

		var command contracts.GenerateImageCommand
		if err := json.Unmarshal(msg.Body, &command); err != nil {
			logger.Error("Failed to unmarshal message", "error", err)
			return err
		}

		event, err := service.GenerateImage(ctx, &command)
		if err != nil {
			logger.Error("Failed to generate image", "error", err)
			return err
		}

		logger.Info("Image generated successfully", "task_id", event.TaskID)
		return nil
	}
}
