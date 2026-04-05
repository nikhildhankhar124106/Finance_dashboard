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
	var dsn string
	// Priority: If a full DATABASE_URL is provided, use it natively.
	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		slog.Info("Using full DATABASE_URL connection string for PostgreSQL.")
	} else {
		// Use individual components with the configured SSL mode.
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		slog.Info("Using decomposed connection parameters for PostgreSQL.", "sslmode", cfg.DBSSLMode)
	}

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

// SyncSequences ensures that the PostgreSQL primary key sequences are in sync with the maximum ID in the tables.
// This prevents "duplicate key value violates unique constraint" errors after manual ID inserts or truncates.
func SyncSequences(db *gorm.DB) error {
	tables := []string{"users", "transactions"}
	for _, table := range tables {
		var maxID int
		// COALESCE(MAX(id), 0) handles empty tables gracefully.
		db.Raw(fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)).Scan(&maxID)
		// Dynamically obtain the sequence name for the table's primary key (usually 'id')
		var sequenceName string
		db.Raw("SELECT pg_get_serial_sequence(?, 'id')", table).Scan(&sequenceName)
		
		if sequenceName == "" {
			// Fallback to standard naming convention if dynamic lookup fails
			sequenceName = fmt.Sprintf("%s_id_seq", table)
		}

		if maxID > 0 {
			slog.Info("Synchronizing sequence", "table", table, "sequence", sequenceName, "max_id", maxID)
			err := db.Exec("SELECT setval(?, ?, true)", sequenceName, maxID).Error
			if err != nil {
				slog.Error("Failed to synchronize sequence", "table", table, "error", err)
				return err
			}
		}
	}
	return nil
}
