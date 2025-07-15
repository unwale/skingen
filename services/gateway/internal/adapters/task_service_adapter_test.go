package adapters

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/unwale/skingen/services/task-service/generated/task/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type mockTaskServer struct {
	pb.UnimplementedTaskServiceServer
}

func (s *mockTaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if req.GetPrompt() == "" {
		return nil, status.Error(codes.InvalidArgument, "prompt is required")
	}

	return &pb.CreateTaskResponse{
		TaskId: 1,
		Status: "pending",
	}, nil
}

func setupTest(t *testing.T) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { lis.Close() }) //nolint:errcheck

	srv := grpc.NewServer()
	pb.RegisterTaskServiceServer(srv, &mockTaskServer{})
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	t.Cleanup(func() { srv.Stop() })

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient(
		"passthrough:///test",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		t.Fatalf("failed to create client connection: %v", err)
	}
	t.Cleanup(func() { conn.Close() }) //nolint:errcheck

	return conn
}

func TestTaskServiceAdapter_CreateTask(t *testing.T) {
	conn := setupTest(t)

	adapter := NewTaskServiceAdapter(conn)

	req := &pb.CreateTaskRequest{Prompt: "My Test Task"}
	resp, err := adapter.CreateTask(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint32(1), resp.TaskId)
}
