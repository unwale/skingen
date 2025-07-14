package domain

import "time"

type Task struct {
	ID        uint
	Prompt    string
	Status    string
	ResultURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}
