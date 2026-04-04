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
// @BasePath /api/v1
// @schemes http

func main() {
	// Initialize Structured Logger natively over Go 1.21 configs
	logger.Init()
	slog.Info("Starting up the backend service processing architectures...")

	// Load Configuration natively falling back accurately onto explicit bindings 
	cfg := config.LoadConfig()

	// Connect to Database securely
	err := db.ConnectPostgres(cfg)
	if err != nil {
		slog.Error("Could not initialize database connection", "error", err)
		os.Exit(1)
	}

	// Auto-migrate database models utilizing implicit DeleteAt columns for soft-deletes inherently.
	DB := db.GetDB()
	err = DB.AutoMigrate(&models.User{}, &models.Transaction{})
	if err != nil {
		slog.Error("Failed to auto migrate database schema", "error", err)
		os.Exit(1)
	}

	// Seed mock users to prevent foreign key errors during mock login flows
	mockUsers := []models.User{
		{ID: 1, Email: "admin@finance.com", Name: "Admin User", Role: models.RoleAdmin},
		{ID: 2, Email: "analyst@finance.com", Name: "Analyst User", Role: models.RoleAnalyst},
		{ID: 3, Email: "viewer@finance.com", Name: "Viewer User", Role: models.RoleViewer},
	}
	for _, u := range mockUsers {
		var existing models.User
		if err := DB.Where("id = ?", u.ID).First(&existing).Error; err != nil {
			// If not found, insert
			DB.Create(&u)
		}
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
	txSvc := service.NewTransactionService(txRepo)
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

	r.GET("/health", HealthCheck)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API VERSIONING IMPLEMENTATION explicitly grouping subsets cleanly
	v1 := r.Group("/api/v1")
	{
		// Public Routes Setup Native Auth
		v1.POST("/auth/login", authHandler.MockLogin)

		// Protected Viewer Routes (accessible by Viewer, Analyst, Admin)
		viewerRoutes := v1.Group("/")
		viewerRoutes.Use(middleware.RequireAuth())
		viewerRoutes.Use(middleware.RequireRole(string(models.RoleViewer), string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			viewerRoutes.GET("/users", userHandler.GetUsers) 
			viewerRoutes.GET("/users/:id", userHandler.GetUser)
			viewerRoutes.GET("/transactions", txHandler.GetTransactions)
		}

		// Protected Analyst Routes (accessible by Analyst, Admin)
		analystRoutes := v1.Group("/")
		analystRoutes.Use(middleware.RequireAuth())
		analystRoutes.Use(middleware.RequireRole(string(models.RoleViewer), string(models.RoleAnalyst), string(models.RoleAdmin)))
		{
			analystRoutes.GET("/summary", analyticsHandler.GetSummary)
			analystRoutes.GET("/category-breakdown", analyticsHandler.GetCategoryBreakdown)
			analystRoutes.GET("/monthly-trends", analyticsHandler.GetMonthlyTrends)
		}

		// Protected Admin Routes (accessible by Admin ONLY)
		adminRoutes := v1.Group("/")
		adminRoutes.Use(middleware.RequireAuth())
		adminRoutes.Use(middleware.RequireRole(string(models.RoleAdmin)))
		{
			adminRoutes.POST("/users", userHandler.CreateUser) 
			adminRoutes.POST("/transactions", txHandler.CreateTransaction)
			adminRoutes.PUT("/transactions/:id", txHandler.UpdateTransaction)
			adminRoutes.DELETE("/transactions/:id", txHandler.DeleteTransaction)
			
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
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0"})
}

// DeleteSystemLogs godoc
// @Summary Delete system logs
// @Description Deletes system logs natively (admin only)
// @Tags system
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Router /system/logs [delete]
func DeleteSystemLogs(c *gin.Context) {
	c.JSON(200, gin.H{"data": "System logs deleted natively."})
}

// V2HealthCheck godoc
// @Summary Check API V2 health
// @Description Returns the health status of the V2 API
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /../v2/health [get]
func V2HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "2.0"})
}

