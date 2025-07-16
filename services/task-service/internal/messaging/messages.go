package messaging

type GenerateImageCommand struct {
	TaskID uint   `json:"task_id"`
	Prompt string `json:"prompt"`
}
