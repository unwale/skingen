package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unwale/skingen/services/gateway/internal/adapters"
	"github.com/unwale/skingen/services/gateway/internal/api/rest"
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

	conn, err := grpc.NewClient(cfg.TaskServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to task service: %v", err)
	}
	defer conn.Close()

	taskServiceAdapter := adapters.NewTaskServiceAdapter(conn)

	service := core.NewGatewayService(taskServiceAdapter)

	httpHandler := rest.NewGatewayHandler(service)

	router := mux.NewRouter()
	httpHandler.RegisterRoutes(router)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
