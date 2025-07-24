package domain

import "time"

const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

type Task struct {
	ID        uint
	Prompt    string
	Status    string
	ObjectID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
