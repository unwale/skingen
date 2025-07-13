package core

type GatewayService interface {
	CreateTask(prompt string) (string, error)
}
