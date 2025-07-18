package core

import (
	"context"

	"github.com/unwale/skingen/pkg/contracts"
)

type workerServiceImpl struct{}

func NewWorkerService() WorkerService {
	return &workerServiceImpl{}
}

func (w *workerServiceImpl) GenerateImage(ctx context.Context, request *contracts.GenerateImageCommand) (*contracts.GenerateImageEvent, error) {
	dummyURL := "http://localhost:0000/image/"

	event := &contracts.GenerateImageEvent{
		TaskID:   request.TaskID,
		ImageURL: dummyURL,
		Status:   contracts.TaskStatusCompleted,
	}

	return event, nil
}
