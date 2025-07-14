package grpc

import (
	"context"

	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"github.com/unwale/skingen/services/task-service/internal/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	pb.UnimplementedTaskServiceServer
	service core.TaskService
}

func NewHandler(service core.TaskService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	task, err := h.service.CreateTask(ctx, req.Prompt)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{
		TaskId:    uint32(task.ID),
		Status:    task.Status,
		CreatedAt: timestamppb.New(task.CreatedAt),
	}, nil
}
