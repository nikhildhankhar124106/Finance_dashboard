package main

import (
	"log"
	"log/slog"
	"time"

	"backend/config"
	"backend/domain/models"
	"backend/infrastructure/db"
	"backend/pkg/logger"
)

func main() {
	logger.Init()
	slog.Info("Starting database seeding process natively bindings limits...")

	// 1. Load Configurations cleanly linking to existing PostgreSQL server definitions
	cfg := config.LoadConfig()

	// 2. Setup connection securely
	if err := db.ConnectPostgres(cfg); err != nil {
		log.Fatalf("Database connection failed during seeding setup: %v", err)
	}

	database := db.GetDB()

	// 3. Clear existing states allowing idempotent execution gracefully 
	database.Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE;")
	database.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")

	slog.Info("Existing table data wiped out securely triggering re-creation mappings...")

	// Automatically map constraints cleanly again
	database.AutoMigrate(&models.User{}, &models.Transaction{})

	// 4. Create Mock Users matching roles identically to our authentication handlers specifications
	admin := models.User{
		Name:  "Admin User",
		Email: "admin@finance.com",
		Role:  models.RoleAdmin,
	}
	analyst := models.User{
		Name:  "Analyst User",
		Email: "analyst@finance.com",
		Role:  models.RoleAnalyst,
	}
	viewer := models.User{
		Name:  "Viewer User",
		Email: "viewer@finance.com",
		Role:  models.RoleViewer,
	}

	database.Create(&admin)
	database.Create(&analyst)
	database.Create(&viewer)
	slog.Info("Mock Users inserted structurally mapping identifiers correctly")

	// 5. Setup basic mock transactions bridging realistically over dates for grouping analytics tests
	now := time.Now()

	transactions := []models.Transaction{
		// Admin assigned transactions
		{
			UserID:   admin.ID,
			Amount:   5000.00,
			Type:     models.TypeIncome,
			Category: "Salary",
			Date:     now.AddDate(0, 0, -5), // 5 Days ago
			Notes:    "Monthly base salary compensation",
		},
		{
			UserID:   admin.ID,
			Amount:   150.00,
			Type:     models.TypeExpense,
			Category: "Groceries",
			Date:     now.AddDate(0, 0, -2),
			Notes:    "Weekly groceries trip",
		},
		
		// Analyst assigned transactions covering multiple historical months enforcing monthly trends logic accurately
		{
			UserID:   analyst.ID,
			Amount:   3200.00,
			Type:     models.TypeIncome,
			Category: "Consulting",
			Date:     now.AddDate(0, -1, -10), // Last Month
			Notes:    "Financial consultation baseline fee",
		},
		{
			UserID:   analyst.ID,
			Amount:   850.50,
			Type:     models.TypeExpense,
			Category: "Rent",
			Date:     now.AddDate(0, -1, -5),
			Notes:    "Property monthly rent",
		},
		{
			UserID:   analyst.ID,
			Amount:   120.00,
			Type:     models.TypeExpense,
			Category: "Utilities",
			Date:     now.AddDate(0, -2, -15), // Two months ago
			Notes:    "Historical electricity and networking bills",
		},

		// Viewer assigned transactions specifically showcasing limited scoping accurately
		{
			UserID:   viewer.ID,
			Amount:   100.00,
			Type:     models.TypeExpense,
			Category: "Entertainment",
			Date:     now,
			Notes:    "Monthly streaming service subscriptions via Viewer",
		},
	}

	for _, tx := range transactions {
		database.Create(&tx)
	}

	slog.Info("Seeded realistic instances spanning history securely natively into PG", "count", len(transactions))
	slog.Info("Database Seeding Finished Successfully!")
}
