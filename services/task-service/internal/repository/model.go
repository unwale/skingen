package repository

import (
	"gorm.io/gorm"

	"github.com/unwale/skingen/services/task-service/internal/domain"
)

type TaskDB struct {
	gorm.Model
	Prompt   string `gorm:"type:text;not null"`
	Status   string `gorm:"type:varchar(20);not null;default:'pending'"`
	ObjectID string `gorm:"type:text"`
}

func (t *TaskDB) TableName() string {
	return "tasks"
}

func (t *TaskDB) toDomain() *domain.Task {
	return &domain.Task{
		ID:        t.ID,
		Prompt:    t.Prompt,
		Status:    t.Status,
		ObjectID:  t.ObjectID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func fromDomain(task *domain.Task) *TaskDB {
	return &TaskDB{
		Model: gorm.Model{
			ID:        task.ID,
			CreatedAt: task.CreatedAt,
			UpdatedAt: task.UpdatedAt,
		},
		Prompt:   task.Prompt,
		Status:   task.Status,
		ObjectID: task.ObjectID,
	}
}
