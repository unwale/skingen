package repository

import (
	"context"

	"github.com/unwale/skingen/services/task-service/internal/core"
	"github.com/unwale/skingen/services/task-service/internal/domain"
	"gorm.io/gorm"
)

type taskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) core.TaskRepository {
	return &taskRepositoryImpl{db: db}
}

func (r *taskRepositoryImpl) SaveTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	taskDB := fromDomain(&task)
	if err := r.db.WithContext(ctx).Create(taskDB).Error; err != nil {
		return domain.Task{}, err
	}
	return *taskDB.toDomain(), nil
}
