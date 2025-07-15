package adapters

import (
	"context"

	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"google.golang.org/grpc"
)

type TaskServiceAdapter struct {
	client pb.TaskServiceClient
}

func NewTaskServiceAdapter(conn *grpc.ClientConn) *TaskServiceAdapter {
	client := pb.NewTaskServiceClient(conn)
	return &TaskServiceAdapter{
		client: client,
	}
}

func (a *TaskServiceAdapter) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	resp, err := a.client.CreateTask(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
