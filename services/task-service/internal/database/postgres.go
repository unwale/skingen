package database

import (
	"fmt"

	"github.com/unwale/skingen/services/task-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
