package main

import (
	"backend/config"
	"backend/infrastructure/db"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	if err := db.ConnectPostgres(cfg); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	database := db.GetDB()
	
	// Check max id
	var maxId int
	database.Raw("SELECT COALESCE(MAX(id), 0) FROM users;").Scan(&maxId)
	fmt.Printf("Current Max ID: %d\n", maxId)

	// Check sequence
	var lastVal int
	database.Raw("SELECT last_value FROM users_id_seq;").Scan(&lastVal)
	fmt.Printf("Sequence last_value: %d\n", lastVal)

	// Fix sequence if needed
	if maxId >= lastVal {
		fmt.Printf("Fixing sequence...\n")
		err := database.Exec(fmt.Sprintf("SELECT setval('users_id_seq', %d);", maxId)).Error
		if err != nil {
			log.Fatalf("Failed to fix sequence: %v", err)
		}
		fmt.Println("Sequence fixed.")
	} else {
		fmt.Println("Sequence is fine.")
	}
}
