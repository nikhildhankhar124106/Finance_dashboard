package db

import (
	"fmt"
	"log/slog"

	"backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres(cfg *config.Config) error {
	// Reverted to dynamic mapping -> Hardcoding proved the Go parser was working accurately, but the DB password itself is rejecting you!
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to postgres explicitly natively", "error", err)
		return err
	}

	DB = db
	slog.Info("Successfully connected to PostgreSQL structured instances natively.")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
