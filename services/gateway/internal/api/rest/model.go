package rest

type CreateTaskRequest struct {
	Prompt string `json:"prompt" validate:"required"`
}

type CreateTaskResponse struct {
	TaskID int `json:"task_id"`
}
