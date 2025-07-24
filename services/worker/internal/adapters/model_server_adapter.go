package adapters

import (
	"context"

	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
	"google.golang.org/grpc"
)

type ModelServerAdapter struct {
	client pb.ImageGeneratorClient
}

func NewTritonAdapter(conn *grpc.ClientConn) *ModelServerAdapter {
	client := pb.NewImageGeneratorClient(conn)
	return &ModelServerAdapter{
		client: client,
	}
}

func (a *ModelServerAdapter) GenerateImage(ctx context.Context, req *pb.GenerateImageRequest) (*pb.GenerateImageResponse, error) {
	resp, err := a.client.GenerateImage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
