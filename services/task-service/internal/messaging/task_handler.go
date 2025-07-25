package messaging

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/unwale/skingen/pkg/contextutil"
	"github.com/unwale/skingen/pkg/contracts"
	"github.com/unwale/skingen/pkg/messaging"
	"github.com/unwale/skingen/services/task-service/internal/core"
)

func CreateTaskResultHandler(service core.TaskService, baseLogger *slog.Logger) messaging.MessageHandler {
	return func(msg amqp091.Delivery) error {
		logger := baseLogger.With("correlation_id", msg.CorrelationId)
		ctx := contextutil.WithLogger(context.Background(), logger)
		ctx = contextutil.WithCorrelationID(ctx, msg.CorrelationId)

		var event contracts.GenerateImageEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		task, err := service.ProcessTaskResult(ctx, event)
		if err != nil {
			log.Printf("Failed to process task result: %v", err)
			return err
		}

		logger.Info("Task result processed successfully", "task_id", task.ID, "status", task.Status)
		return nil
	}
}
