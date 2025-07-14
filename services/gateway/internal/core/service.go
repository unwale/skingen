package core

type gatewayServiceImpl struct {
}

func NewGatewayService() GatewayService {
	return &gatewayServiceImpl{}
}

func (s *gatewayServiceImpl) CreateTask(prompt string) (string, error) {
	return "dummy-task-id", nil
}
