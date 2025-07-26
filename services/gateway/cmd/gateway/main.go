package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unwale/skingen/pkg/logging"
	"github.com/unwale/skingen/services/gateway/internal/adapters"
	"github.com/unwale/skingen/services/gateway/internal/api/grpc/interceptors"
	"github.com/unwale/skingen/services/gateway/internal/api/rest"
	"github.com/unwale/skingen/services/gateway/internal/api/rest/middleware"
	"github.com/unwale/skingen/services/gateway/internal/config"
	"github.com/unwale/skingen/services/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logging.NewLogger(cfg.ServiceName, cfg.LoggingLevel)
	logger.Info("Starting gateway service", "port", cfg.Port)

	conn, err := grpc.NewClient(
		cfg.TaskServiceUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.CorrelationIDInterceptor()),
	)
	if err != nil {
		logger.Error("Failed to connect to task service", "error", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("Failed to close gRPC connection", "error", err)
		}
	}()

	taskServiceAdapter := adapters.NewTaskServiceAdapter(conn)

	service := core.NewGatewayService(taskServiceAdapter, logger)

	httpHandler := rest.NewGatewayHandler(service)

	router := mux.NewRouter()
	httpHandler.RegisterRoutes(router)

	loggingMiddleware := middleware.LoggingMiddleware(logger)
	router.Use(loggingMiddleware)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Error("Failed to start HTTP server", "error", err)
	}
}
