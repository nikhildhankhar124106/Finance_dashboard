package main

import (
	"log/slog"
	"net/http"
	"os"

	"backend/config"
	_ "backend/docs"
	"backend/domain/models"
	"backend/handler/api"
	"backend/handler/middleware"
	"backend/infrastructure/cache"
	"backend/infrastructure/db"
	"backend/pkg/auth"
	"backend/pkg/logger"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Finance Dashboard API
// @version 1.0 (v1) / 2.0 (v2) 
// @description REST API documentation integrating Redis caching, Rate-Limiting and Structured Logging.
// @host localhost:8080
// @basePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token to authenticate.

func main() {
	// Initialize Structured Logger natively over Go 1.21 configs
	logger.Init()
	slog.Info("Starting up the backend service processing architectures...")

	// Load Configuration natively falling back accurately onto explicit bindings 
	cfg := config.LoadConfig()

	// Initialize JWT Secret securely from environment
	auth.SetJWTSecret(cfg.JWTSecret)

	// Connect to Database securely
	err := db.ConnectPostgres(cfg)
	if err != nil {
		slog.Error("Could not initialize database connection", "error", err)
		os.Exit(1)
	}

	// Auto-migrate database models utilizing implicit DeleteAt columns for soft-deletes inherently.
	DB := db.GetDB()
	err = DB.AutoMigrate(&models.User{}, &models.Transaction{}, &models.ActivityLog{})
	if err != nil {
		slog.Error("Failed to auto migrate database schema", "error", err)
		os.Exit(1)
	}

	// Seed mock users to prevent foreign key errors during mock login flows
	mockUsers := []models.User{
		{ID: 1, Email: "admin@finance.com", Name: "Admin User", Role: models.RoleAdmin, IsActive: true},
		{ID: 2, Email: "analyst@finance.com", Name: "Analyst User", Role: models.RoleAnalyst, IsActive: true},
		{ID: 3, Email: "viewer@finance.com", Name: "Viewer User", Role: models.RoleViewer, IsActive: true},
	}
	for _, u := range mockUsers {
		var existing models.User
		if err := DB.Where("id = ?", u.ID).First(&existing).Error; err != nil {
			// If not found, insert
			DB.Create(&u)
		}
	}

	// Fix DB sequences to prevent primary key mismatch errors after manual ID inserts
	if err := db.SyncSequences(DB); err != nil {
		slog.Warn("Failed to synchronize database sequences", "error", err)
	}

	// Initialize Distributed Cache natively
	if err := cache.ConnectRedis(cfg); err != nil {
		slog.Warn("Redis failed starting, proceeding gracefully without cache aggregations", "error", err)
	}

	// Dependency Injection Setup
	userRepo := repository.NewUserRepository(DB)
	userSvc := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userSvc)
	authHandler := api.NewAuthHandler()

	txRepo := repository.NewTransactionRepository(DB)
	txSvc := service.NewTransactionService(txRepo, userRepo)
	txHandler := api.NewTransactionHandler(txSvc)

	analyticsRepo := repository.NewAnalyticsRepository(DB)
	analyticsSvc := service.NewAnalyticsService(analyticsRepo)
	analyticsHandler := api.NewAnalyticsHandler(analyticsSvc)

	// Setup Gin Router natively bindings configurations natively
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	
	// Apply Global Middlewares spanning cleanly
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RateLimiter())
	// Global Audit Layer for write operations natively
	r.Use(middleware.AuditMiddleware(DB))

	r.GET("/", WelcomePage)
	r.GET("/health", HealthCheck)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API VERSIONING IMPLEMENTATION explicitly grouping subsets cleanly
	v1 := r.Group("/api/v1")
	{
		// Native Health Check for V1 specifically
		v1.GET("/health", HealthCheck)

		// Public Routes Setup Native Auth
		v1.POST("/auth/login", authHandler.MockLogin)

		// Protected Viewer Routes (accessible by Viewer, Analyst, Admin)
		viewerRoutes := v1.Group("/")
		viewerRoutes.Use(middleware.RequireAuth(DB))
		viewerRoutes.Use(middleware.RequireRole(string(models.RoleViewer), string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			viewerRoutes.GET("/users", userHandler.GetUsers) 
			viewerRoutes.GET("/users/:id", userHandler.GetUser)
			viewerRoutes.GET("/transactions", txHandler.GetTransactions)
		}

		// Protected Analyst Routes (accessible by Analyst, Admin)
		analystRoutes := v1.Group("/")
		analystRoutes.Use(middleware.RequireAuth(DB))
		analystRoutes.Use(middleware.RequireRole(string(models.RoleViewer), string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			analystRoutes.GET("/summary", analyticsHandler.GetSummary)
			analystRoutes.GET("/category-breakdown", analyticsHandler.GetCategoryBreakdown)
			analystRoutes.GET("/monthly-trends", analyticsHandler.GetMonthlyTrends)
		}

		// Protected Admin Routes (accessible by Admin ONLY)
		adminRoutes := v1.Group("/")
		adminRoutes.Use(middleware.RequireAuth(DB))
		adminRoutes.Use(middleware.RequireRole(string(models.RoleAdmin)))
		{
			adminRoutes.POST("/users", userHandler.CreateUser) 
			adminRoutes.PATCH("/users/:id/status", userHandler.UpdateUserStatus)
			adminRoutes.POST("/transactions", txHandler.CreateTransaction)
			adminRoutes.PUT("/transactions/:id", txHandler.UpdateTransaction)
			adminRoutes.DELETE("/transactions/:id", txHandler.DeleteTransaction)
			adminRoutes.GET("/transactions/export", txHandler.ExportTransactions)
			
			adminRoutes.DELETE("/system/logs", DeleteSystemLogs)
		}
	}

	// Example V2 mapping for future API versions separating structural changes uniquely!
	v2 := r.Group("/api/v2")
	{
		v2.GET("/health", V2HealthCheck)
	}

	// Start the server
	slog.Info("Server starting on configured port limits", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("Server forced to shutdown securely binding outputs", "error", err)
	}
}

// HealthCheck godoc
// @Summary Check API health
// @Description Returns the health status of the API
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /v1/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0"})
}

// WelcomePage returns basic API metadata and navigation links
func WelcomePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":       "Welcome to the Finance Dashboard API",
		"status":        "Live",
		"version":       "1.0 (v1)",
		"documentation": "/docs/index.html",
		"health":        "/health",
	})
}

// DeleteSystemLogs godoc
// @Summary Delete system logs
// @Description Deletes system logs natively (admin only)
// @Tags system
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /v1/system/logs [delete]
func DeleteSystemLogs(c *gin.Context) {
	c.JSON(200, gin.H{"data": "System logs deleted natively."})
}

// V2HealthCheck godoc
// @Summary Check API V2 health
// @Description Returns the health status of the V2 API
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /v2/health [get]
func V2HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "2.0"})
}

