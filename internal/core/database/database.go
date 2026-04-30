package database

import (
	"fmt"
	"jsunnykhan/go-clean-template/internal/core/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(cfg config.DatabaseConfig, models ...any) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormCfg)

	if err != nil {
		return nil, fmt.Errorf("postgres: open connection: %w", err)

	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("postgres: get sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return nil, fmt.Errorf("postgres: enable uuid-ossp extension: %w", err)
	}

	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return nil, fmt.Errorf("postgres: auto migrate: %w", err)
		}
	}

	return db, nil
}
