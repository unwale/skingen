package rest

type CreateTaskRequest struct {
	Prompt string `json:"prompt" validate:"required"`
}

type CreateTaskResponse struct {
	TaskID string `json:"task_id"`
}
