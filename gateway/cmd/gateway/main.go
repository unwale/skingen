package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unwale/skingen/gateway/internal/api/rest"
	"github.com/unwale/skingen/gateway/internal/config"
	"github.com/unwale/skingen/gateway/internal/core"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service := core.NewGatewayService()

	httpHandler := rest.NewGatewayHandler(service)

	router := mux.NewRouter()
	httpHandler.RegisterRoutes(router)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
