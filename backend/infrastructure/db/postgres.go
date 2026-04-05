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

// SyncSequences ensures that the PostgreSQL primary key sequences are in sync with the maximum ID in the tables.
// This prevents "duplicate key value violates unique constraint" errors after manual ID inserts or truncates.
func SyncSequences(db *gorm.DB) error {
	tables := []string{"users", "transactions"}
	for _, table := range tables {
		var maxID int
		// COALESCE(MAX(id), 0) handles empty tables gracefully.
		db.Raw(fmt.Sprintf("SELECT COALESCE(MAX(id), 0) FROM %s", table)).Scan(&maxID)
		
		// setval sets the sequence's current value. The next call to nextval will return maxID + 1.
		// If maxID is 0, we don't need to do anything as the sequence starts at 1 by default.
		if maxID > 0 {
			slog.Info("Synchronizing sequence for table", "table", table, "max_id", maxID)
			err := db.Exec(fmt.Sprintf("SELECT setval('%s_id_seq', %d)", table, maxID)).Error
			if err != nil {
				slog.Error("Failed to synchronize sequence", "table", table, "error", err)
				return err
			}
		}
	}
	return nil
}
