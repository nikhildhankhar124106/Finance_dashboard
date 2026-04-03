package db

import (
	"fmt"
	"backend/config"
	"backend/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error.Printf("Failed to connect to postgres: %v", err)
		return err
	}

	DB = db
	logger.Info.Println("Successfully connected to PostgreSQL")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
