package database

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{}
	if cfg.AppEnv == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("PostgreSQL connection successfully established")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database auto-migration...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Comment{},
		&model.Like{},
		&model.Follow{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}
