package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/unwale/skingen/services/gateway/internal/core"
)

type gatewayHandler struct {
	service core.GatewayService
}

func NewGatewayHandler(service core.GatewayService) *gatewayHandler {
	return &gatewayHandler{service: service}
}

func (h *gatewayHandler) RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/create-task", h.CreateTaskHandler).Methods("POST")
}

func (h *gatewayHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validator := validator.New(validator.WithRequiredStructEnabled())
	if err := validator.Struct(req); err != nil {

		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := h.service.CreateTask(r.Context(), req.Prompt)
	if err != nil {
		http.Error(w, "Failed to create task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateTaskResponse{TaskID: taskID}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
