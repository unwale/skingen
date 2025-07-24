package core

import (
	"context"

	"github.com/unwale/skingen/pkg/contracts"
	pb "github.com/unwale/skingen/services/model-server/generated/model/v1"
)

type WorkerService interface {
	GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error)
}

type ModelServer interface {
	GenerateImage(ctx context.Context, request *pb.GenerateImageRequest) (*pb.GenerateImageResponse, error)
}
