package core

import (
	"context"

	"github.com/unwale/skingen/pkg/contracts"
)

type WorkerService interface {
	GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error)
}
