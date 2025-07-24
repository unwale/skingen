package contracts

const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

type GenerateImageCommand struct {
	TaskID uint   `json:"task_id"`
	Prompt string `json:"prompt"`
}

type GenerateImageEvent struct {
	TaskID   uint   `json:"task_id"`
	ObjectID string `json:"object_id"`
	Status   string `json:"status"`
}
