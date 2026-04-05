package main

import (
	"backend/config"
	"backend/infrastructure/db"
	"fmt"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	if err := db.ConnectPostgres(cfg); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	database := db.GetDB()
	var maxId int
	database.Raw("SELECT COALESCE(MAX(id), 0) FROM users;").Scan(&maxId)

	var lastVal int
	database.Raw("SELECT last_value FROM users_id_seq;").Scan(&lastVal)

	f, _ := os.Create("tmp_out_2.txt")
	defer f.Close()
	fmt.Fprintf(f, "Max ID in users: %d\n", maxId)
	fmt.Fprintf(f, "Sequence last_value: %d\n", lastVal)
}
