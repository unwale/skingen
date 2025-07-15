package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGatewayService struct {
	mock.Mock
}

func (m *MockGatewayService) CreateTask(ctx context.Context, prompt string) (int, error) {
	args := m.Called(ctx, prompt)
	return args.Int(0), args.Error(1)
}

func TestCreateTaskHandler(t *testing.T) {
	mockGatewayService := new(MockGatewayService)
	handler := NewGatewayHandler(mockGatewayService)

	req := "{\"prompt\": \"test\"}"
	mockGatewayService.On("CreateTask", mock.Anything, "test").Return(1, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/create-task", strings.NewReader(req))
	handler.CreateTaskHandler(w, r)

	mockGatewayService.AssertExpectations(t)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	var resp CreateTaskResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.TaskID)

}
