package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/unwale/skingen/services/task-service/internal/domain"
	"gorm.io/gorm"
)

func TestToDomain(t *testing.T) {
	taskDB := &TaskDB{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Prompt:    "Test Prompt",
		Status:    "pending",
		ResultUrl: "http://example.com/result",
	}

	task := taskDB.toDomain()

	assert.Equal(t, taskDB.ID, task.ID)
	assert.Equal(t, taskDB.Prompt, task.Prompt)
	assert.Equal(t, taskDB.Status, task.Status)
	assert.Equal(t, taskDB.ResultUrl, task.ResultURL)
	assert.Equal(t, taskDB.CreatedAt, task.CreatedAt)
	assert.Equal(t, taskDB.UpdatedAt, task.UpdatedAt)
}

func TestFromDomain(t *testing.T) {
	task := &domain.Task{
		ID:        1,
		Prompt:    "Test Prompt",
		Status:    "pending",
		ResultURL: "http://example.com/result",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	taskDB := fromDomain(task)

	assert.Equal(t, task.ID, taskDB.ID)
	assert.Equal(t, task.Prompt, taskDB.Prompt)
	assert.Equal(t, task.Status, taskDB.Status)
	assert.Equal(t, task.ResultURL, taskDB.ResultUrl)
	assert.Equal(t, task.CreatedAt, taskDB.CreatedAt)
	assert.Equal(t, task.UpdatedAt, taskDB.UpdatedAt)
}
